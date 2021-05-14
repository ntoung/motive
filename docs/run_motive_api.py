import os
import sys
import time
from collections import defaultdict

import pandas as pd
import requests
from tqdm import tqdm

pd.options.mode.chained_assignment = None  # default='warn'

# Fill in API key here
MOTIVE_API_URL = "https://api-data.motivesoftware.com/scores"
MOTIVE_API_KEY = "ADD-API-KEY-HERE"


def batch_apply_func_to_df(df, func, batch_size=250):
    """ Applies the given function to the dataframe
        in batches of batch_size rows, returns
        concatenated dataframe of results
    """
    in_dfs = [df[i:i + batch_size] for i in range(0, df.shape[0], batch_size)]
    out_dfs = []
    for batch_df in tqdm(in_dfs):
        out_dfs.append(func(batch_df))
    df = pd.concat(out_dfs)
    return df


def score_with_motive_api(df, models=None, domain='cx', channel='surveys'):
    """ Score a dataframe with the Motive API.
    Expects text to be scored to be in 'text' column in dataframe.
    Unless models are specified, default to returning sentiment and emotion.
    Domain and channel parameters are also defaulted for ease of use.
    """
    if not models:
        models = ['sentiment', 'emotion']
    if isinstance(models, str):
        models = [models]
    if df.shape[0] > 1000:
        print("Max batch size is 1000, truncating")
        df = df.head(n=1000)

    # Format documents appropriately
    docs_to_score = [{'document_id': str(i), 'text': s} for i, s in enumerate(df['text'].values)]

    # Set up scoring call
    headers = {"X-API-Key": MOTIVE_API_KEY}
    payload = {
        "correlation_id": "0",
        "domain": domain,
        "data_channel": channel,
        "models": models,
        "documents": docs_to_score
    }
    response = requests.post(MOTIVE_API_URL, json=payload, headers=headers)

    # Poll for response
    if response.status_code == requests.codes.accepted:
        job_id = response.json()['job_id']

        # get results
        results_url = f"{MOTIVE_API_URL}/{job_id}"
        response = requests.get(results_url, headers=headers)

        # get responses
        while response.status_code == requests.codes.accepted \
                and response.json()['status'] not in ['ERROR', 'DONE']:
            time.sleep(1)
            response = requests.get(results_url, headers=headers)
        all_responses = response.json()
    else:
        print(f'Error {response.content.decode("utf-8")}')

    # Get the top classification for each model,
    # ignoring secondary classifications (eg, emotions 2 & 3)
    # or sentence-level scores for simplicity:
    classifications = defaultdict(dict)
    for doc in all_responses["documents"]:
        idx = doc['document_id']
        for model in doc["models"]:
            model_name = model['model']
            if not model['scores']:
                continue
            top_label = max(model['scores'], key=lambda x: x['score'])
            label, score = top_label['label'], top_label['score']
            classifications[idx][model_name] = label
            classifications[idx]['%s_score' % model_name] = score

    # Add classifications to original dataframe:
    doc_ids = [d['document_id'] for d in docs_to_score]
    for model in models:
        df[model] = [classifications.get(idx, {}).get(model, None) for idx in doc_ids]
        score_col = '%s_score' % model
        df[score_col] = [classifications.get(idx, {}).get(score_col, None) for idx in doc_ids]

    # Return scored dataframe
    return df


if __name__ == '__main__':

    # Pass in filename to score via command-line:
    # (Should we upgrade to named parameters?)
    if len(sys.argv) == 2:
        file_name = sys.argv[1]
        if not os.path.exists(file_name):
            print("Error: Unrecognized file '%s'" % file_name)
            sys.exit(1)
    else:
        print(f"You need to call {sys.argv[0]} with the name of the file to process (csv, xlsx)")
        sys.exit(1)

    # Check to make sure this is a CSV of an XLSX
    if not file_name.endswith(".xlsx") and not file_name.endswith(".csv"):
        print("Error: Expecting CSV of XLSX file")
        sys.exit(1)
    filetype = '.csv' if file_name.endswith('.csv') else '.xlsx'

    # Create scored file name by appending
    # "_Motive" to the filename:
    scored_file_name = file_name.replace(filetype, '_Motive' + filetype)

    # Assume that text to score is in 'text' column,
    # or specify a different column here:
    text_column = "text"

    # Read the file
    if filetype == '.csv':
        df = pd.read_csv(file_name)
    else:
        df = pd.read_excel(file_name)
    original_cols = list(df.columns)

    # Check for a text column:
    if text_column not in df.columns:
        print("Error: Expecting a 'text' column with documents to score")
        sys.exit(1)

    # Fill empty text rows with empty string
    df[text_column] = df[text_column].fillna("")

    # Report on total number of rows to user
    print()
    print(len(df), "rows to be processed")

    # Score in batches of 1000
    df = batch_apply_func_to_df(df, score_with_motive_api, batch_size=1000)

    # Write results to CSV of XLSX based on input format:
    if filetype == '.csv':
        df.to_csv(scored_file_name, index=False)
    else:
        df.to_excel(scored_file_name, index=False, engine='xlsxwriter')

    # Identify the columns we've added so we can call them out to the user:
    additional_cols = [c for c in df.columns if c not in original_cols]
    print("Added columns:", additional_cols)

    print("Scored file written to '%s'" % scored_file_name)
    print("Done.")
    print()
