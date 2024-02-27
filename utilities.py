import pandas as pd
import sqlite3

def json_to_dataframe(json_data):
    """
    Convert JSON data to a pandas DataFrame.
    
    Parameters:
        json_data (dict): JSON data to convert.
        
    Returns:
        DataFrame: Pandas DataFrame created from the JSON data.
    """
    df = pd.DataFrame(json_data)

    return df



def read_table(table_name):
    """
    Read data from the specified table in the SQLite database.
    
    Parameters:
        table_name (str): The name of the table to read from.
    
    Returns:
        pd.DataFrame: A DataFrame containing the data from the specified table.
    """
    try:
        # Connect to SQLite database and read data into DataFrame
        conn = sqlite3.connect('crypto_data.db')
        df = pd.read_sql_query(f"SELECT * FROM {table_name}", conn)
        
        # Close the connection
        conn.close()
        
        return df
    except sqlite3.Error as e:
        print("Error reading data from the database:", e)
        return pd.DataFrame()  # Return an empty DataFrame in case of error

    

def append_df_to_sql(df, table_name, database_name='crypto_data.db'):
    """
    Append a DataFrame to an existing table in a SQLite database if the table exists.
    If the table does not exist, create the table and write the DataFrame to it.

    Parameters:
        df (pd.DataFrame): The DataFrame to append to the table.
        table_name (str): The name of the table in the SQLite database.
        database_name (str): The name of the SQLite database file.
    """
    try:
        # Connect to SQLite database
        conn = sqlite3.connect(database_name)

        # Append DataFrame to existing SQLite table if it exists, create table if it doesn't exist
        df.to_sql(table_name, conn, if_exists='append', index=False)

        # Commit the transaction
        conn.commit()

        # Close the connection
        conn.close()

        print(f"DataFrame has been successfully appended to the SQLite database table '{table_name}'.")
    except Exception as e:
        print("Error:", e)

# Example usage:
# append_df_to_sql(df, 'your_table_name', 'your_database.db')
