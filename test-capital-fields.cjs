#!/usr/bin/env node

// Simple end-to-end test script to verify the capital fields integration
const http = require('http');

const BASE_URL = 'http://localhost:8080';

function makeRequest(method, path, data = null) {
    return new Promise((resolve, reject) => {
        const options = {
            hostname: 'localhost',
            port: 8080,
            path: path,
            method: method,
            headers: {
                'Content-Type': 'application/json',
            }
        };

        const req = http.request(options, (res) => {
            let body = '';
            res.on('data', (chunk) => {
                body += chunk;
            });
            res.on('end', () => {
                try {
                    const jsonBody = JSON.parse(body);
                    resolve({ status: res.statusCode, data: jsonBody });
                } catch (e) {
                    resolve({ status: res.statusCode, data: body });
                }
            });
        });

        req.on('error', (err) => {
            reject(err);
        });

        if (data) {
            req.write(JSON.stringify(data));
        }
        
        req.end();
    });
}

async function runTests() {
    console.log('🚀 Starting end-to-end capital fields test...\n');
    
    const timestamp = Date.now();
    const testUsername = `e2euser_${timestamp}`;
    const testBankName = `E2E Bank ${timestamp}`;
    
    try {
        // Test 1: Create new user
        console.log('1️⃣ Creating new user...');
        const createResult = await makeRequest('POST', '/api/user', {
            username: testUsername,
            bankName: testBankName
        });
        
        if (createResult.status !== 200) {
            throw new Error(`User creation failed: ${createResult.status} - ${JSON.stringify(createResult.data)}`);
        }
        
        console.log('✅ User created successfully');
        console.log(`   Username: ${createResult.data.username}`);
        console.log(`   Bank Name: ${createResult.data.bankName}`);
        console.log(`   Claimed Capital: £${createResult.data.claimedCapital}`);
        console.log(`   Actual Capital: £${createResult.data.actualCapital}`);
        
        if (createResult.data.claimedCapital !== 1000) {
            throw new Error(`Expected claimed capital 1000, got ${createResult.data.claimedCapital}`);
        }
        if (createResult.data.actualCapital !== 1000) {
            throw new Error(`Expected actual capital 1000, got ${createResult.data.actualCapital}`);
        }
        
        // Test 2: Login with created user
        console.log('\n2️⃣ Logging in with created user...');
        const loginResult = await makeRequest('POST', '/api/login', {
            username: testUsername
        });
        
        if (loginResult.status !== 200) {
            throw new Error(`Login failed: ${loginResult.status} - ${JSON.stringify(loginResult.data)}`);
        }
        
        console.log('✅ Login successful');
        console.log(`   Username: ${loginResult.data.username}`);
        console.log(`   Bank Name: ${loginResult.data.bankName}`);
        console.log(`   Claimed Capital: £${loginResult.data.claimedCapital}`);
        console.log(`   Actual Capital: £${loginResult.data.actualCapital}`);
        
        if (loginResult.data.claimedCapital !== 1000) {
            throw new Error(`Expected claimed capital 1000, got ${loginResult.data.claimedCapital}`);
        }
        if (loginResult.data.actualCapital !== 1000) {
            throw new Error(`Expected actual capital 1000, got ${loginResult.data.actualCapital}`);
        }
        
        console.log('\n🎉 All tests passed! Capital fields integration is working correctly.');
        console.log('\n📋 Summary:');
        console.log('   ✅ New users get initialized with £1000 claimed and actual capital');
        console.log('   ✅ Capital values are properly returned in API responses');
        console.log('   ✅ Login endpoint returns complete user data including capital fields');
        
    } catch (error) {
        console.error(`\n❌ Test failed: ${error.message}`);
        process.exit(1);
    }
}

// Check if server is running first
console.log('🔍 Checking if backend server is running...');
makeRequest('GET', '/').then(() => {
    runTests();
}).catch(() => {
    console.error('❌ Backend server is not running on http://localhost:8080');
    console.log('💡 Please start the backend server first:');
    console.log('   cd backend && go run main.go');
    process.exit(1);
});
