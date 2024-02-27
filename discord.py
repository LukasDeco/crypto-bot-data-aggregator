import requests
import pandas as pd

def get_discord_info(discord_url):

    api_url = f'{discord_url}?with_counts=true'
    response = requests.get(api_url)
    
    print(api_url)
    
    # Check if the request was successful (status code 200)
    if response.status_code == 200:
        # Parse the JSON response
        data = response.json()
        return data
    else:
        print("Error:", response.status_code)
        return None
    
data = {
    'MATIC': 'https://discord.com/api/invite/XvpHAxZ',
    'VET': 'https://discord.com/api/invite/vechain',
    'ETH': 'https://discord.com/api/invite/ethereum-org',
    'DIMO': 'https://discord.com/api/invite/dimonetwork',
    'MXC': 'https://discord.com/api/invite/mxcfoundation'
}

# Create DataFrame from dictionary
df = pd.DataFrame(list(data.items()), columns=['Token', 'Discord_URL'])

print(df)