import http from 'k6/http';
import { check, sleep } from 'k6';

// Stress test configuration
export const options = {
  stages: [
    { duration: '1m', target: 50 },   // Ramp up to 50 users
    { duration: '2m', target: 100 }, // Ramp up to 100 users
    { duration: '2m', target: 200 }, // Ramp up to 200 users
    { duration: '2m', target: 400 }, // Ramp up to 400 users
    { duration: '2m', target: 600 }, // Ramp up to 600 users
    { duration: '2m', target: 800 }, // Ramp up to 800 users
    { duration: '2m', target: 1000 }, // Ramp up to 1000 users
    { duration: '5m', target: 1000 }, // Stay at 1000 users
    { duration: '2m', target: 0 },   // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<2000'], // 95% of requests must complete below 2s
    http_req_failed: ['rate<0.2'],     // Error rate must be below 20%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Stress test Product Service with concurrent requests
  const productPromises = [];
  for (let i = 0; i < 5; i++) {
    productPromises.push(http.get(`${BASE_URL}/api/products`));
  }
  
  const productResponses = productPromises;
  productResponses.forEach(response => {
    check(response, {
      'product service stress test': (r) => r.status === 200,
    });
  });

  // Stress test Basket Service
  const userId = `stress_user_${Math.floor(Math.random() * 10000)}`;
  
  // Create basket
  const basketResponse = http.get(`${BASE_URL}/api/baskets/${userId}`);
  check(basketResponse, {
    'basket creation stress test': (r) => r.status === 200,
  });

  // Add multiple items rapidly
  if (basketResponse.status === 200) {
    for (let i = 0; i < 10; i++) {
      const addItemPayload = JSON.stringify({
        product_id: Math.floor(Math.random() * 100),
        quantity: Math.floor(Math.random() * 5) + 1,
        name: `Stress Test Product ${i}`,
        price: Math.random() * 100,
        category: 'Stress Test'
      });
      
      const addItemResponse = http.post(`${BASE_URL}/api/baskets/${userId}/items`, addItemPayload, {
        headers: { 'Content-Type': 'application/json' },
      });
      
      check(addItemResponse, {
        'add item stress test': (r) => r.status === 200,
      });
    }
  }

  // Stress test Payment Service
  const paymentPayload = JSON.stringify({
    user_id: userId,
    amount: Math.random() * 1000,
    currency: 'USD',
    method: 'credit_card',
    provider: 'stripe'
  });
  
  const paymentResponse = http.post(`${BASE_URL}/api/payments`, paymentPayload, {
    headers: { 'Content-Type': 'application/json' },
  });
  
  check(paymentResponse, {
    'payment stress test': (r) => r.status === 200,
  });

  // Stress test Health Endpoints
  const healthPromises = [];
  for (let i = 0; i < 3; i++) {
    healthPromises.push(http.get(`${BASE_URL}/health`));
  }
  
  const healthResponses = healthPromises;
  healthResponses.forEach(response => {
    check(response, {
      'health stress test': (r) => r.status === 200,
    });
  });

  // Random sleep to simulate real user behavior
  sleep(Math.random() * 2);
}

export function handleSummary(data) {
  return {
    'stress-test-results.json': JSON.stringify(data, null, 2),
  };
}
