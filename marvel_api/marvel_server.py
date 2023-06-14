from flask import Flask, render_template
import requests
import os

app = Flask(__name__)

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/characters')
def get_characters():
    api_key = os.environ.get('MARVEL_API_KEY')  # Obtiene la clave de API de las variables de entorno
    url = f'https://gateway.marvel.com/v1/public/characters?apikey={api_key}'
    response = requests.get(url)
    data = response.json()
    characters = data['data']['results']
    return render_template('characters.html', characters=characters)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
