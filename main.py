import reddit
import github
from dotenv import load_dotenv
import os
import subprocess



def main():

    load_dotenv()
    github_token = os.getenv('github_token')

    reddit.get_reddit_main()
    github.get_github_main(github_token)
    
    
    # executing golang data collection process 
    executable_path = './go-exec'  # Use the appropriate path to your executable

    # Running the executable and waiting for it to complete
    result = subprocess.run([executable_path], capture_output=True, text=True)

    # Printing the stdout and stderr of the executable
    print("STDOUT:", result.stdout)
    print("STDERR:", result.stderr)




if __name__ == '__main__':
    main()
