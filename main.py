import reddit
import github
from dotenv import load_dotenv
import os



def main():

    load_dotenv()
    github_token = os.getenv('github_token')

    reddit.get_reddit_main()
    github.get_github_main(github_token)




if __name__ == '__main__':
    main()
