import os
import pandas as pd

def merge_csv_files(input_folder, output_file):
    # Define the range of files to be merged (PP1 to PP13)
    file_range = range(1, 14)  # Files from PP1 to PP13
    
    # Create an empty DataFrame to store merged data
    merged_df = pd.DataFrame()
    
    # Iterate over the file range and read each CSV into a DataFrame
    for i in file_range:
        # Construct the filename dynamically for each file in the range
        filename = os.path.join(input_folder, f"New_York_PP{i}.csv_Building.csv")
        
        # Check if the file exists before attempting to read
        if os.path.exists(filename):
            print(f"Merging {filename}...")
            df = pd.read_csv(filename)
            merged_df = pd.concat([merged_df, df], ignore_index=True)
        else:
            print(f"File {filename} does not exist. Skipping.")
    
    # Save the merged DataFrame into a new CSV file
    merged_df.to_csv(output_file, index=False)
    print(f"Merged CSV saved as {output_file}")

# Specify the input folder and output file name
input_folder = "path_to_your_csv_files"  # Change this to the folder where your CSV files are located
output_file = "merged_output.csv"  # Change this to your desired output file name

# Call the function to merge the CSV files
merge_csv_files(input_folder, output_file)
