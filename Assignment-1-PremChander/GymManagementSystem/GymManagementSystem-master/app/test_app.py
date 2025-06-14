import pytest
from flask import Flask
from app import app as flask_app

@pytest.fixture
def client():
    flask_app.config['TESTING'] = True
    with flask_app.test_client() as client:
        yield client

def test_home_page(client):
    rv = client.get('/')
    assert rv.status_code == 200

def test_login_get(client):
    rv = client.get('/login')
    assert rv.status_code == 200

def test_404(client):
    rv = client.get('/nonexistent')
    assert rv.status_code == 404
