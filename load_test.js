import http from 'k6/http';
import { sleep } from 'k6';
import { check } from 'k6';

const BASE_URL = 'http:localhost:8080'
const USERS = 100; 
let jwtTokens = Array(USERS).fill(null); 
let balances = Array(USERS).fill(1000); 


function authenticate(userIndex) {
    const authRequest = {
        username: `testuser${userIndex}`, 
        password: 'password',
    };

    const res = http.post(`${BASE_URL}/api/auth`, JSON.stringify(authRequest), {
        headers: { 'Content-Type': 'application/json' },
    });

    check(res, {
        'auth status is 200': (r) => r.status === 200,
    });

    const responseBody = JSON.parse(res.body);
    if (responseBody.token) {
        balances[userIndex] = 1000; 
    }
    
    return responseBody.token; 
}

export const options = {
    vus: USERS, 
    duration: '30s', 
};

export default function () {
    const userIndex = __VU - 1; 

    if (!jwtTokens[userIndex]) {
        jwtTokens[userIndex] = authenticate(userIndex); 
    }

    
    let res1 = http.get(`${BASE_URL}/api/info`, {
        headers: { 'Authorization': `Bearer ${jwtTokens[userIndex]}` },
    });
    check(res1, {
        'info status is 200': (r) => r.status === 200,
    });

    
    if (balances[userIndex] > 0) {
        let sendCoinRequest = {
            toUser: `testuser${(userIndex + 1) % USERS}`, 
            amount: 10, 
        };
        
        let res2 = http.post(`${BASE_URL}/api/sendCoin`, JSON.stringify(sendCoinRequest), {
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${jwtTokens[userIndex]}`,
            },
        });
        check(res2, {
            'sendCoin status is 200': (r) => r.status === 200,
        });

        balances[userIndex] -= 10; 
    }

    
    if (balances[userIndex] > 0) {
        let itemToBuy = 'book'; 
        let res3 = http.get(`${BASE_URL}/api/buy/${itemToBuy}`, {
            headers: { 'Authorization': `Bearer ${jwtTokens[userIndex]}` },
        });
        check(res3, {
            'buy item status is 200': (r) => r.status === 200,
        });

        balances[userIndex] -= 50; 
    }

    sleep(1); 
}