const baseUrl = 'http://localhost:8080/branch';

// Create query parameters
const params = new URLSearchParams({
  key1: 'value1',
  key2: 'value2',
}).toString();

fetch(`${baseUrl}?${params}`)
  .then(response => {
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json(); // or response.text() if your endpoint returns plain text
  })
  .then(data => {
    console.log('GET Response:', data);
  })
  .catch(error => {
    console.error('GET Error:', error);
  });



const url = 'http://localhost:8080/branch';

const jsonData = {
  name: 'branch 1',
};

fetch(url, {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(jsonData),
})
  .then(response => {
    if (!response.ok) throw new Error('Network response was not ok');
    return response.json(); // or response.text() if your endpoint returns plain text
  })
  .then(data => {
    console.log('POST Response:', data);
  })
  .catch(error => {
    console.error('POST Error:', error);
  });
