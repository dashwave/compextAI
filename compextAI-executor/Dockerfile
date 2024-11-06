FROM python:3.10-slim

# Set the working directory
WORKDIR /app

# Copy requirements.txt into the container
COPY requirements.txt .

# Create a virtual environment
RUN python -m venv venv

# Activate the virtual environment
RUN . venv/bin/activate

# Install dependencies from requirements.txt using the virtual environment's pip
RUN pip install -r requirements.txt

# Copy the rest of the application code
COPY . .

EXPOSE 8889
# Command to run the FastAPI server
CMD ["python", "api.py"]
