from HttpExample import ___init__ as function

class MockRequest:
    def __init__(self):
        pass

def test_main():
    req = MockRequest()
    res = function.main(req)
    assert res.status_code == 200
    assert res.get_body().decode() == "Hello, world!"
