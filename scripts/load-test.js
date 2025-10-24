import http from 'k6/http';
import { check, sleep } from 'k6';

// Test configuration
export const options = {
  stages: [
    { duration: '2m', target: 100 }, // Ramp up to 100 users
    { duration: '5m', target: 100 }, // Stay at 100 users
    { duration: '2m', target: 200 }, // Ramp up to 200 users
    { duration: '5m', target: 200 }, // Stay at 200 users
    { duration: '2m', target: 0 },   // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
    http_req_failed: ['rate<0.1'],    // Error rate must be below 10%
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  // Test Product Service
  const productResponse = http.get(`${BASE_URL}/api/products`);
  check(productResponse, {
    'product service status is 200': (r) => r.status === 200,
    'product service response time < 500ms': (r) => r.timings.duration < 500,
  });

  // Test specific product
  if (productResponse.status === 200) {
    const products = JSON.parse(productResponse.body);
    if (products.length > 0) {
      const productId = products[0].id;
      const productDetailResponse = http.get(`${BASE_URL}/api/products/${productId}`);
      check(productDetailResponse, {
        'product detail status is 200': (r) => r.status === 200,
        'product detail response time < 300ms': (r) => r.timings.duration < 300,
      });
    }
  }

  // Test Basket Service
  const userId = `user_${Math.floor(Math.random() * 1000)}`;
  const basketResponse = http.get(`${BASE_URL}/api/baskets/${userId}`);
  check(basketResponse, {
    'basket service status is 200': (r) => r.status === 200,
    'basket service response time < 400ms': (r) => r.timings.duration < 400,
  });

  // Test adding item to basket
  if (basketResponse.status === 200) {
    const addItemPayload = JSON.stringify({
      product_id: 1,
      quantity: 2,
      name: 'Test Product',
      price: 29.99,
      category: 'Electronics'
    });
    
    const addItemResponse = http.post(`${BASE_URL}/api/baskets/${userId}/items`, addItemPayload, {
      headers: { 'Content-Type': 'application/json' },
    });
    
    check(addItemResponse, {
      'add item status is 200': (r) => r.status === 200,
      'add item response time < 300ms': (r) => r.timings.duration < 300,
    });
  }

  // Test Payment Service
  const paymentPayload = JSON.stringify({
    user_id: userId,
    amount: 59.98,
    currency: 'USD',
    method: 'credit_card',
    provider: 'stripe'
  });
  
  const paymentResponse = http.post(`${BASE_URL}/api/payments`, paymentPayload, {
    headers: { 'Content-Type': 'application/json' },
  });
  
  check(paymentResponse, {
    'payment service status is 200': (r) => r.status === 200,
    'payment service response time < 1000ms': (r) => r.timings.duration < 1000,
  });

  // Test Health Endpoints
  const healthResponse = http.get(`${BASE_URL}/health`);
  check(healthResponse, {
    'health check status is 200': (r) => r.status === 200,
    'health check response time < 100ms': (r) => r.timings.duration < 100,
  });

  // Test Metrics Endpoint
  const metricsResponse = http.get(`${BASE_URL}/metrics`);
  check(metricsResponse, {
    'metrics status is 200': (r) => r.status === 200,
    'metrics response time < 200ms': (r) => r.timings.duration < 200,
  });

  sleep(1);
}

export function handleSummary(data) {
  return {
    'load-test-results.json': JSON.stringify(data, null, 2),
  };
}
