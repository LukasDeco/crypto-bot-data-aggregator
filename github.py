import pandas as pd
import requests
import time
import utilities as util

def is_valid_org(org_name, github_token):
    url = f"https://api.github.com/orgs/{org_name}"
    
    headers = {'Authorization': f'Bearer {github_token}'}
    
    print(url)
    response = requests.get(url, headers=headers)
    return response.status_code == 200

def get_all_repos(org_name, github_token):
    if not is_valid_org(org_name, github_token):
        print(f"The organization '{org_name}' is not a valid GitHub organization.")
        return []
    
    repos = []
    page = 1
    headers = {'Authorization': f'Bearer {github_token}'}
    while True:
        headers = {'Authorization': f'token {github_token}'} if github_token else {}
        repos_url = f"https://api.github.com/orgs/{org_name}/repos?per_page=100&page={page}"
        response = requests.get(repos_url, headers=headers)
        if response.status_code == 200:
            page_repos = response.json()
            if not page_repos:
                break
            repos.extend(page_repos)
            page += 1
        else:
            print(f"Failed to fetch repositories for {org_name}. Status code: {response.status_code}")
            break
    return repos

# Define a function to normalize and extract repository data
def extract_repo_data(repo, token):
    dict = {
        'token': token,
        'id': repo['id'],
        'node_id': repo['node_id'],
        'name': repo['name'],
        'full_name': repo['full_name'],
        'private': repo['private'],
        'owner_login': repo['owner']['login'],
        'owner_id': repo['owner']['id'],
        'owner_type': repo['owner']['type'],
        'html_url': repo['html_url'],
        'description': repo['description'],
        'fork': repo['fork'],
        'url': repo['url'],
        'created_at': repo['created_at'],
        'updated_at': repo['updated_at'],
        'pushed_at': repo['pushed_at'],
        'git_url': repo['git_url'],
        'ssh_url': repo['ssh_url'],
        'clone_url': repo['clone_url'],
        'svn_url': repo['svn_url'],
        'homepage': repo['homepage'],
        'size': repo['size'],
        'stargazers_count': repo['stargazers_count'],
        'watchers_count': repo['watchers_count'],
        'language': repo['language'],
        'has_issues': repo['has_issues'],
        'has_projects': repo['has_projects'],
        'has_downloads': repo['has_downloads'],
        'has_wiki': repo['has_wiki'],
        'has_pages': repo['has_pages'],
        'has_discussions': repo.get('has_discussions', False),  # Use .get for keys that might not exist in all repos
        'forks_count': repo['forks_count'],
        'archived': repo['archived'],
        'disabled': repo['disabled'],
        'open_issues_count': repo['open_issues_count'],
        'license_key': repo['license']['key'] if repo['license'] else None,
        'license_name': repo['license']['name'] if repo['license'] else None,
        'allow_forking': repo['allow_forking'],
        'is_template': repo['is_template'],
        'topics': ', '.join(repo['topics']),
        'visibility': repo['visibility'],
        'forks': repo['forks'],
        'open_issues': repo['open_issues'],
        'watchers': repo['watchers'],
        'default_branch': repo['default_branch'],
    }
    
    df = pd.DataFrame([dict])
    
    return df

def get_github_main(github_token):
    df = pd.read_csv(r"input_data/github_repositories.csv")
    for index, row in df.iterrows():
        repo_response = get_all_repos(row['org_name'], github_token)
        time.sleep(5)
        for repo in repo_response:
            github_df = extract_repo_data(repo, row['token'])
            util.append_df_to_sql(github_df, "github")
            
            