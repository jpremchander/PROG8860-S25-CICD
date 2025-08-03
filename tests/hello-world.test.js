const { helloWorldHandler } = require('../HelloWorld/index');

// Mock the Azure Functions context
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
    body
});

describe('Azure Function Hello World Tests', () => {
    let mockContext;

    beforeEach(() => {
        mockContext = createMockContext();
    });

    // Test Case 1: Basic HTTP response returns "Hello, World!"
    test('should return Hello, World! for default request', async () => {
        const mockRequest = createMockRequest();
        
        const response = await helloWorldHandler(mockRequest, mockContext);

        expect(response.status).toBe(200);
        expect(response.body).toBe('Hello, World!');
        expect(response.headers['Content-Type']).toBe('text/plain');
    });

    // Test Case 2: HTTP response code is 200
    test('should return 200 status code', async () => {
        const mockRequest = createMockRequest();
        
        const response = await helloWorldHandler(mockRequest, mockContext);

        expect(response.status).toBe(200);
    });

    // Test Case 3: Edge case - custom name parameter from query
    test('should return personalized greeting with custom name from query', async () => {
        const queryParams = new Map([['name', 'Jenkins']]);
        const mockRequest = createMockRequest('http://localhost:7071/api/HelloWorld?name=Jenkins', queryParams);
        
        const response = await helloWorldHandler(mockRequest, mockContext);

        expect(response.status).toBe(200);
        expect(response.body).toBe('Hello, Jenkins!');
        expect(response.headers['Content-Type']).toBe('text/plain');
    });

    // Test Case 4: Edge case - name from POST body
    test('should return personalized greeting from POST body', async () => {
        const mockRequest = createMockRequest('http://localhost:7071/api/HelloWorld', new Map(), { name: 'Azure' });
        
        const response = await helloWorldHandler(mockRequest, mockContext);

        expect(response.status).toBe(200);
        expect(response.body).toBe('Hello, Azure!');
        expect(response.headers['Content-Type']).toBe('text/plain');
    });

    // Test Case 5: Edge case - empty name parameter defaults to World
    test('should return Hello, World! for empty name', async () => {
        const queryParams = new Map([['name', '']]);
        const mockRequest = createMockRequest('http://localhost:7071/api/HelloWorld?name=', queryParams);
        
        const response = await helloWorldHandler(mockRequest, mockContext);

        expect(response.status).toBe(200);
        expect(response.body).toBe('Hello, World!');
    });

    // Test Case 6: Verify context logging is called
    test('should log the request URL', async () => {
        const mockRequest = createMockRequest();
        
        await helloWorldHandler(mockRequest, mockContext);

        expect(mockContext.log).toHaveBeenCalledWith(
            expect.stringContaining('Http function processed request for url')
        );
    });
});
