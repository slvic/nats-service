async function getOrder() {
    let orderId = document.getElementById('msgId').value;
    const response = await fetch(`http://localhost:3000/messages?id=${orderId}`);
    const result = await response.json()
    document.getElementById('message').textContent = JSON.stringify(result, undefined, 2)
}
