import requests



def get_community_info(community_name):

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

