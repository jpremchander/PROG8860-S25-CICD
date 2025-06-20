# YelpCamp Go - Complete Endpoint Validation Guide

## 🔍 **How to Test All Endpoints**

### **1. Web Interface Endpoints (Browser)**

#### **Homepage**
\`\`\`
URL: http://localhost:3000/
Expected: Beautiful landing page with hero section and working images
\`\`\`

#### **All Campgrounds**
\`\`\`
URL: http://localhost:3000/campgrounds
Expected: Grid of campground cards with images and details
\`\`\`

#### **Individual Campground**
\`\`\`
URL: http://localhost:3000/campgrounds/1
Expected: Detailed campground page with amenities
\`\`\`

#### **Registration Page**
\`\`\`
URL: http://localhost:3000/register
Expected: Registration form with validation
Test: Fill form and submit - should show success message
\`\`\`

#### **Login Page**
\`\`\`
URL: http://localhost:3000/login
Expected: Login form with demo credentials
Test Credentials:
- Username: demo
- Password: password
\`\`\`

---

### **2. API Endpoints (JSON) - Use curl or Postman**

#### **Health Check**
\`\`\`bash
curl http://localhost:3000/api/health
\`\`\`
**Expected Response:**
\`\`\`json
{
  "status": "ok",
  "message": "YelpCamp Go API is running",
  "timestamp": "2024-06-19",
  "version": "1.0.0"
}
\`\`\`

#### **Get All Campgrounds**
\`\`\`bash
curl http://localhost:3000/api/campgrounds
\`\`\`
**Expected Response:**
\`\`\`json
{
  "message": "Campgrounds retrieved successfully",
  "data": [
    {
      "id": 1,
      "title": "Sunset Valley Campground",
      "description": "A beautiful campground with stunning sunset views",
      "location": "Yosemite National Park, CA",
      "price": 35.99,
      "image": "https://images.unsplash.com/..."
    }
  ],
  "count": 3
}
\`\`\`

#### **Get Single Campground**
\`\`\`bash
curl http://localhost:3000/api/campgrounds/1
\`\`\`
**Expected Response:**
\`\`\`json
{
  "message": "Campground retrieved successfully",
  "campground": {
    "id": "1",
    "title": "Sunset Valley Campground",
    "description": "A beautiful campground with stunning sunset views over the valley",
    "location": "Yosemite National Park, CA",
    "price": 35.99,
    "image": "https://images.unsplash.com/...",
    "amenities": ["Fire Pits", "Restrooms", "Showers", "Picnic Tables", "Hiking Trails"]
  }
}
\`\`\`

#### **User Registration**
\`\`\`bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
\`\`\`
**Expected Response:**
\`\`\`json
{
  "message": "User registered successfully",
  "user": {
    "id": 123,
    "username": "testuser",
    "email": "test@example.com"
  },
  "token": "sample_jwt_token_here"
}
\`\`\`

#### **User Login**
\`\`\`bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "demo",
    "password": "password"
  }'
\`\`\`
**Expected Response:**
\`\`\`json
{
  "message": "Login successful",
  "user": {
    "id": 123,
    "username": "demo",
    "email": "demo@yelpcamp.com"
  },
  "token": "sample_jwt_token_here"
}
\`\`\`

#### **Test Endpoint**
\`\`\`bash
curl http://localhost:3000/api/test
\`\`\`
**Expected Response:**
\`\`\`json
{
  "message": "Test endpoint working",
  "data": {
    "server_time": "2024-06-19T10:00:00Z",
    "endpoints": [
      "/api/health",
      "/api/campgrounds",
      "/api/campgrounds/:id",
      "/api/auth/register",
      "/api/auth/login"
    ]
  }
}
\`\`\`

---

### **3. Error Testing**

#### **Invalid Registration Data**
\`\`\`bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "",
    "email": "invalid-email",
    "password": "123"
  }'
\`\`\`
**Expected Response:**
\`\`\`json
{
  "error": "Validation failed",
  "details": "Key: 'username' Error:Field validation for 'username' failed on the 'required' tag"
}
\`\`\`

#### **Invalid Login Credentials**
\`\`\`bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "wronguser",
    "password": "wrongpass"
  }'
\`\`\`
**Expected Response:**
\`\`\`json
{
  "error": "Invalid credentials"
}
\`\`\`

---

### **4. Complete Validation Checklist**

#### **✅ Web Interface Tests**
- [ ] Homepage loads with images
- [ ] Campgrounds page shows all cards with images
- [ ] Individual campground page displays correctly
- [ ] Registration form works and shows success message
- [ ] Login form works with demo credentials (demo/password)
- [ ] Navigation between pages works
- [ ] Responsive design works on mobile

#### **✅ API Tests**
- [ ] `/api/health` returns status ok
- [ ] `/api/campgrounds` returns array of campgrounds
- [ ] `/api/campgrounds/1` returns single campground
- [ ] `/api/auth/register` accepts valid user data
- [ ] `/api/auth/login` works with demo credentials
- [ ] `/api/test` returns endpoint list
- [ ] Error handling works for invalid data

#### **✅ Image Loading Tests**
- [ ] Hero section background image loads
- [ ] Homepage camping image loads
- [ ] All campground card images load
- [ ] Individual campground detail image loads

---

### **5. Quick Test Script**

Save this as `test-endpoints.sh`:

\`\`\`bash
#!/bin/bash
echo "🧪 Testing YelpCamp Go Endpoints..."

echo "1. Health Check:"
curl -s http://localhost:3000/api/health | jq .

echo -e "\n2. Campgrounds:"
curl -s http://localhost:3000/api/campgrounds | jq '.message, .count'

echo -e "\n3. Single Campground:"
curl -s http://localhost:3000/api/campgrounds/1 | jq '.campground.title'

echo -e "\n4. Registration Test:"
curl -s -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","password":"password123"}' | jq '.message'

echo -e "\n5. Login Test:"
curl -s -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"password"}' | jq '.message'

echo -e "\n✅ All tests completed!"
\`\`\`

Run with: `chmod +x test-endpoints.sh && ./test-endpoints.sh`

---

## 🎯 **Expected Results After Fixes:**

1. **✅ Images Load:** All images now use Unsplash URLs that work
2. **✅ Registration Works:** Form submits to API with proper validation
3. **✅ Login Works:** Demo credentials (demo/password) work
4. **✅ All Endpoints Respond:** Complete API functionality
5. **✅ Error Handling:** Proper error messages for invalid data
6. **✅ Professional UI:** Bootstrap styling with working forms

Your YelpCamp application should now be fully functional! 🎉
