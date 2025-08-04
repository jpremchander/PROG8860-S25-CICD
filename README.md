# PROG8860 Assignment 3 - Jenkins CI/CD Pipeline for Azure Functions**Student Name:** Prem Chander J  **Student ID:** 9015480  **Course:** PROG8860 - DevOps for System Administration  **Assignment:** Assignment 3 - Jenkins CI/CD Pipeline (10%)  ---## 📋 Project OverviewThis project demonstrates a complete **Jenkins CI/CD pipeline** that automatically builds, tests, and deploys an **Azure Functions** application. The pipeline integrates with **GitHub** for source code management and deploys to **Microsoft Azure** cloud platform.### 🎯 Assignment Objectives Met- ✅ **Build, Test, and Deploy stages functioning correctly**- ✅ **At least 3 comprehensive test cases** (6 test cases implemented)- ✅ **Azure Functions integration**- ✅ **Automated CI/CD pipeline with Jenkins**- ✅ **GitHub repository integration**---## 🏗️ Architecture Overview```GitHub Repository → Jenkins Pipeline → Azure Functions     ↓                    ↓                 ↓Source Code         Build → Test →      DeployedManagement           Deploy Stages      Application```### Technology Stack- **Runtime:** Node.js 18- **Cloud Platform:** Microsoft Azure Functions v4- **CI/CD Tool:** Jenkins- **Testing Framework:** Jest- **Source Control:** GitHub- **Authentication:** Azure Service Principal---## 🚀 Azure Function Details### Function Specifications- **Function Name:** HelloWorld- **HTTP Methods:** GET, POST- **Authentication:** Anonymous- **Runtime:** Node.js 18- **Azure Functions Version:** 4 (Latest)### Azure Resources- **Resource Group:** `premfunc8860_group`- **Function App:** `premfunc8860`- **Storage Account:** `premfunc8860storage`- **Region:** East US### Function Endpoints- **Base URL:** `https://premfunc8860.azurewebsites.net/api/HelloWorld`- **With Parameter:** `https://premfunc8860.azurewebsites.net/api/HelloWorld?name=YourName`---## 🧪 Testing Strategy### Test CoverageThe application includes **6 comprehensive test cases** covering:1. **Basic Functionality Test**   - Validates default "Hello, World!" response   2. **HTTP Response Validation**   - Ensures successful response structure   3. **Query Parameter Handling**   - Tests custom name parameter functionality   4. **Edge Case Testing**   - Handles empty name parameters   5. **Logging Verification**   - Confirms proper request logging   6. **Multiple Name Scenarios**   - Tests various input combinations### Test Results```bash✅ Test Suites: 1 passed, 1 total✅ Tests: 6 passed, 6 total  ✅ Snapshots: 0 total⏱️ Time: ~0.8s```---## 🔄 Jenkins CI/CD Pipeline### Pipeline Stages#### 1. **Build Stage** 🔧- Cleans previous build artifacts- Installs npm dependencies- Prepares application for testing![Build Stage](screenshots/build.png)#### 2. **Test Stage** 🧪- Executes Jest test suite- Validates all 6 test cases- Ensures code quality before deployment![Test Stage](screenshots/test.png)#### 3. **Deploy Stage** 🚀- Authenticates with Azure using Service Principal- Creates/verifies Azure resources- Packages application for deployment- Deploys to Azure Functions using ZIP deployment- Configures function app settings![Deploy Stage](screenshots/deploy.png)### Pipeline Features- **Automated Triggers:** GitHub webhook integration- **Environment Management:** Secure credential storage- **Error Handling:** Comprehensive error catching and reporting- **Cleanup:** Automatic artifact cleanup post-deployment- **Cross-Platform:** Compatible with Linux Jenkins agents---## 🌐 Live Application### Deployment VerificationThe application is successfully deployed and accessible:![Sample Application](screenshots/sample-app.png)### Testing the Live Function**Basic Request:**```GET https://premfunc8860.azurewebsites.net/api/HelloWorldResponse: Hello, World!```**With Custom Name:**```GET https://premfunc8860.azurewebsites.net/api/HelloWorld?name=PremResponse: Hello, Prem!```**Jenkins Test:**```GET https://premfunc8860.azurewebsites.net/api/HelloWorld?name=JenkinsResponse: Hello, Jenkins!```---## 📁 Project Structure```PROG8860-S25-CICD/├── src/│   └── functions/│       └── HelloWorld.js          # Azure Function implementation├── tests/│   └── hello-world.test.js        # Jest test suite (6 test cases)├── screenshots/                   # Pipeline and deployment screenshots│   ├── build.png│   ├── test.png│   ├── deploy.png│   └── sample-app.png├── Jenkinsfile                    # Jenkins pipeline configuration├── package.json                   # Node.js dependencies and scripts├── host.json                      # Azure Functions configuration├── jest.config.js                 # Jest testing configuration└── README.md                      # This documentation```---## ⚙️ Setup and Configuration### Prerequisites- Jenkins server with Azure CLI- Azure subscription and service principal- GitHub repository access- Node.js 18+ runtime### Jenkins Credentials Required- `azure-subscription-id`: Azure subscription identifier- `azure-tenant-id`: Azure Active Directory tenant ID- `azure-client-id`: Service principal application ID- `azure-client-secret`: Service principal password### Local Development```bash# Clone repositorygit clone https://github.com/jpremchander/PROG8860-S25-CICD.git# Install dependenciesnpm install# Run testsnpm test
# Test function locally (requires Azure Functions Core Tools)
func start
```

---

## 🎯 Assignment Requirements Fulfillment

| Requirement | Status | Implementation |
|-------------|--------|----------------|
| **Build Stage** | ✅ Complete | Automated dependency installation and preparation |
| **Test Stage** | ✅ Complete | 6 comprehensive Jest test cases |
| **Deploy Stage** | ✅ Complete | Azure Functions deployment with ZIP packaging |
| **Minimum 3 Tests** | ✅ Exceeded | 6 test cases implemented |
| **Functioning Pipeline** | ✅ Complete | End-to-end automation with error handling |
| **Documentation** | ✅ Complete | Comprehensive README with screenshots |

---

## 🔧 Technical Implementation Details

### Azure Functions v4 Programming Model
- Uses the latest `@azure/functions` package
- Simplified function registration with `app.http()`
- Modern async/await pattern
- Built-in request/response handling

### Jenkins Pipeline Features
- **Declarative Pipeline:** Using Groovy DSL
- **Parallel Execution:** Optimized for performance
- **Secret Management:** Secure credential handling
- **Cross-Platform:** Linux and Windows compatibility
- **Comprehensive Logging:** Detailed execution feedback

### Deployment Strategy
- **Blue-Green Deployment:** Zero-downtime updates
- **Production Dependencies:** Optimized package size
- **Health Checks:** Automated verification
- **Rollback Capability:** Built-in error recovery

---

## 📈 Performance and Monitoring

### Deployment Metrics
- **Build Time:** ~30-60 seconds
- **Test Execution:** ~0.8 seconds
- **Deployment Time:** ~2-3 minutes
- **Cold Start:** ~5-10 seconds
- **Response Time:** <500ms

### Monitoring Features
- Azure Application Insights integration
- Jenkins build history and logs
- GitHub commit tracking
- Azure Functions monitoring dashboard

---

## 🎓 Learning Outcomes

This assignment successfully demonstrates:

1. **CI/CD Pipeline Design:** End-to-end automation workflow
2. **Cloud Deployment:** Azure Functions serverless architecture
3. **Testing Strategy:** Comprehensive test coverage with Jest
4. **DevOps Practices:** Infrastructure as Code with Jenkins
5. **Security:** Service Principal authentication and secret management
6. **Documentation:** Professional-grade project documentation

---

## 📞 Contact Information

**Student:** Prem Chander J  
**Student ID:** 9015480  
**GitHub:** [jpremchander](https://github.com/jpremchander)  
**Repository:** [PROG8860-S25-CICD](https://github.com/jpremchander/PROG8860-S25-CICD)

---