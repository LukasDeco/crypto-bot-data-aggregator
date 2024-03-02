import requests
from datetime import datetime
import time
import pandas as pd
import utilities as util

def format_n_save_reddit(json_data, subreddit, token):
    df = pd.DataFrame(json_data)
    transposed_df = df.transpose()

    df = transposed_df.reset_index()

    current_datetime = datetime.now()
    df['datetime'] = current_datetime
    
    df['subreddit'] = subreddit
    df['token'] = token

    df = df[df['index'] == 'data']
    df = df[['datetime', 'subreddit', 'token', 'accounts_active', 'active_user_count', 'subscribers']]
    
    return df


def get_subreddit_info(community_name):

    api_url = f'https://www.reddit.com/r/{community_name}/about.json'
    headers = {'User-Agent': 'YourBot/0.2'}  # Add a user-agent header to avoid being blocked
    response = requests.get(api_url, headers=headers)
    
    # Check if the request was successful (status code 200)
    if response.status_code == 200:
        # Parse the JSON response
        data = response.json()
        return data
    else:
        print("Error:", response.status_code)
        return None

def get_reddit_main():
    # Read in Reddit CSV
    df = pd.read_csv(r"input_data/reddit_communities.csv")
    
    # Loop through Rows
    for index, row in df.iterrows():
        print(f'Start - {row["token"]}')
        json_response = get_subreddit_info(row['subreddit'])
        subreddit_df = format_n_save_reddit(json_response, row['subreddit'], row['token'])
        util.append_df_to_sql(subreddit_df, "reddit")
        #time.sleep(1)