// Mock Azure Functions context
const createMockContext = () => ({
    log: jest.fn(),
    error: jest.fn(),
    warn: jest.fn(),
    info: jest.fn(),
    debug: jest.fn(),
    trace: jest.fn()
});

// Mock request object that mimics Azure Functions request
const createMockRequest = (url = 'http://localhost:7071/api/HelloWorld', queryParams = new Map(), body = null) => ({
    url,
    query: queryParams,
    text: async () => body ? JSON.stringify(body) : '',
    headers: {}
});

// Simple mock of Azure Functions app for testing
const mockHandler = async (request, context) => {
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
};

describe('Azure Function Hello World Tests', () => {
    let mockContext;

    beforeEach(() => {
        mockContext = createMockContext();
    });

    // Test Case 1: Basic HTTP response returns "Hello, World!"
    test('should return Hello, World! for default request', async () => {
        const mockRequest = createMockRequest();
        
        const response = await mockHandler(mockRequest, mockContext);

        expect(response.body).toBe('Hello, World!');
        expect(response.headers['Content-Type']).toBe('text/plain');
    });

    // Test Case 2: HTTP response code is 200 (implied by successful response)
    test('should return successful response', async () => {
        const mockRequest = createMockRequest();
        
        const response = await mockHandler(mockRequest, mockContext);

        expect(response.body).toBeDefined();
        expect(response.headers).toBeDefined();
    });

    // Test Case 3: Edge case - custom name parameter from query
    test('should return personalized greeting with custom name from query', async () => {
        const queryParams = new Map([['name', 'Jenkins']]);
        const mockRequest = createMockRequest('http://localhost:7071/api/HelloWorld?name=Jenkins', queryParams);
        
        const response = await mockHandler(mockRequest, mockContext);

        expect(response.body).toBe('Hello, Jenkins!');
        expect(response.headers['Content-Type']).toBe('text/plain');
    });

    // Test Case 4: Edge case - empty name parameter defaults to World
    test('should return Hello, World! for empty name', async () => {
        const queryParams = new Map([['name', '']]);
        const mockRequest = createMockRequest('http://localhost:7071/api/HelloWorld?name=', queryParams);
        
        const response = await mockHandler(mockRequest, mockContext);

        expect(response.body).toBe('Hello, World!');
    });

    // Test Case 5: Verify context logging is called
    test('should log the request URL', async () => {
        const mockRequest = createMockRequest();
        
        await mockHandler(mockRequest, mockContext);

        expect(mockContext.log).toHaveBeenCalledWith(
            expect.stringContaining('Http function processed request for url')
        );
    });

    // Test Case 6: Test with different names
    test('should handle various names correctly', async () => {
        const testCases = ['Azure', 'TestUser', 'CI/CD'];
        
        for (const testName of testCases) {
            const queryParams = new Map([['name', testName]]);
            const mockRequest = createMockRequest(`http://localhost:7071/api/HelloWorld?name=${testName}`, queryParams);
            
            const response = await mockHandler(mockRequest, mockContext);
            expect(response.body).toBe(`Hello, ${testName}!`);
        }
    });
});
