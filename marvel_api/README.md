Certainly! Here's an example of how the `README.md` file for your small server that consumes the Marvel API could look in English:

```markdown
# Marvel API Server

This is a basic HTTP server that consumes the Marvel API and displays the results on a simple web page.

## Requirements

- Python 3.9 or higher installed on your system.
- Marvel API key. You can obtain it by registering on the Marvel Developer website.

## Installation

1. Clone this repository to your local machine:

   ```bash
   git clone [<REPOSITORY_URL>](https://github.com/nicaurybenitez/k8s_tooling.git)
   ```

2. Navigate to the project directory:

   ```bash
   cd marvel-api-server
   ```

3. Create and activate a virtual environment (optional but recommended):

   ```bash
   python3 -m venv venv
   source venv/bin/activate
   ```

4. Install the project dependencies:

   ```bash
   pip install -r requirements.txt
   ```

## Configuration

1. Obtain a Marvel API key by registering on the Marvel Developer website.

2. Open the `docker-compose.yml` file and replace `YOUR_API_KEY` on the line `- MARVEL_API_KEY=YOUR_API_KEY` with your own Marvel API key.

## Usage

1. Run the server:

   ```bash
   python marvel_server.py
   ```

2. Open your web browser and visit `http://localhost:5000` to access the main page.

3. Click on the "View Characters" link to see the Marvel characters retrieved from the API.

## Contribution

If you wish to contribute to this project, feel free to do so. You can open an issue to report bugs or request new features.

