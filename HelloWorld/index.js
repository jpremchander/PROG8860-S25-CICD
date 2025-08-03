const { app } = require('@azure/functions');

// Main handler function - export this for testing
const helloWorldHandler = async (request, context) => {
    context.log(`Http function processed request for url "${request.url}"`);

    const name = request.query.get('name') || request.body?.name || 'World';
    
    return { 
        status: 200,
        body: `Hello, ${name}!`,
        headers: {
            'Content-Type': 'text/plain'
        }
    };
};

app.http('HelloWorld', {
    methods: ['GET', 'POST'],
    authLevel: 'anonymous',
    handler: helloWorldHandler
});

// Export both the app and handler for testing
module.exports = { app, helloWorldHandler };
