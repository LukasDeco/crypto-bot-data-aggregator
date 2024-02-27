import reddit
import utilities as util
from datetime import datetime
import time

def format_n_save_reddit(json_data, community):
    df = util.json_to_dataframe(json_data)
    transposed_df = df.transpose()

    df = transposed_df.reset_index()

    current_datetime = datetime.now()
    df['datetime'] = current_datetime

    df = df[df['index'] == 'data']
    df = df[['datetime', 'accounts_active', 'active_user_count', 'subscribers']]
    util.append_df_to_sql(df, community)

def main():
    reddit_communities = util.read_table('reddit_communities')

    for community in reddit_communities['Reddit_Community']:
        json_response = reddit.get_community_info(community)
        format_n_save_reddit(json_response, community)
        time.sleep(5)




if __name__ == '__main__':
    main()
