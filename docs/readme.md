Motive API Starter Code Notes

To get started, first install the required dependencies, as defined in the included requirements.txt file. Ideally, use a new virtual environment using the tooling of your choice in order to isolate things appropriately. Note also that this code was tested on Python 3.7. 
pip install -r requirements.txt

Next, open up run_motive_api.py.  

You will want to copy in your assigned MOTIVE_API_KEY at the top of the file.
The code assumes that your input file has a column named "text". If that's not the case, please either rename the column in your input file, or modify the following line in the code to match your field name:

text_column = "text"

Finally, you are ready to run!
python run_motive_api.py <file_name>

Assuming your file is in the same folder as the python script, this should write a scored output file into the same directory, and it will keep the same extension and file format as your original file (CSV or XLSX). 
