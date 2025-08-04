const { app } = require('@azure/functions');

app.http('HelloWorld', {
    methods: ['GET', 'POST'],
    authLevel: 'anonymous',
    handler: async (request, context) => {
        context.log(`Http function processed request for url "${request.url}"`);

        // Get name from query parameter, request body, or default to 'World'
        const name = request.query.get('name') || 
                    (await request.text().catch(() => '')) || 
                    'World';

        return { 
            body: `Hello, ${name}!`,
            headers: {
                'Content-Type': 'text/plain'
            }
        };
    }
});

// Export for testing
module.exports = app;
