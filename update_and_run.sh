#!/bin/bash

# Navigate to your repository directory
cd /home/steven/crypto-bot-data-aggregator

# Set the executable bit for the script
chmod +x update_and_run.sh

# Pull the latest changes from the main branch
git pull origin main

# If your project depends on Python packages listed in a requirements.txt file,
# uncomment the next line to update them. This requires pip to be installed.
pip3 install -r requirements.txt

# Replace "your_script.py" with the actual name of your Python script
python3 main.py

# Add the database file to the staging area (replace 'your_database.db' with your actual database file name)
git add crypto_data.db

# Commit the changes (customize the commit message as needed)
git commit -m "Update database file with new data"

# Push the changes back to GitHub
git push origin main

echo "Database file committed and pushed to GitHub successfully."