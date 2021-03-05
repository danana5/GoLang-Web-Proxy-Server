class Server:
    def __init__ (self, config):
        self.config = config

    def test_print(self):
        print(self.config)

test = Server("testing...")
test.test_print()

    
