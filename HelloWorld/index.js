const { app } = require('@azure/functions');

// Main handler function - export this for testing
const helloWorldHandler = async (request, context) => {
    context.log(`Http function processed request for url "${request.url}"`);

    try {
        // Handle query parameter or body name, with fallback to 'World'
        let name = 'World';
        
        // Try to get name from query parameters first
        if (request.query && request.query.get) {
            const queryName = request.query.get('name');
            if (queryName && queryName.trim()) {
                name = queryName.trim();
            }
        }
        
        // If no valid query name, try body
        if (name === 'World' && request.body && request.body.name && request.body.name.trim()) {
            name = request.body.name.trim();
        }

        return { 
            status: 200,
            body: `Hello, ${name}!`,
            headers: {
                'Content-Type': 'text/plain'
            }
        };
    } catch (error) {
        context.log.error('Error processing request:', error);
        return {
            status: 500,
            body: 'Internal Server Error',
            headers: {
                'Content-Type': 'text/plain'
            }
        };
    }
};

app.http('HelloWorld', {
    methods: ['GET', 'POST'],
    authLevel: 'anonymous',
    handler: helloWorldHandler
});

// Export both the app and handler for testing
module.exports = { app, helloWorldHandler };
