#!/bin/bash

# Pull the latest changes from the main branch
git pull origin main

# If your project depends on Python packages listed in a requirements.txt file,
# uncomment the next line to update them. This requires pip to be installed.
pip3 install -r requirements.txt

# Replace "your_script.py" with the actual name of your Python script
python3 main.py
