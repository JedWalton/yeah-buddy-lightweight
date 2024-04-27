<script>
// @ts-nocheck

    import { browser } from '$app/env';
    
    let username = '';
    let password = '';

    async function loginUser() {
        if (!browser) return; // Ensures this code only runs on the client side

        const backendUrl = `http://localhost:8080`;

        const response = await fetch(`${backendUrl}/auth/public/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username, password })
        });

        const data = await response.json();
        if (response.ok) {
            alert('Login successful');
            console.log('Token:', data);
        } else {
            alert('Login failed: ' + data.error);
        }
    }
</script>

<div>
    <h2>Login</h2>
    <input type="text" bind:value={username} placeholder="Username">
    <input type="password" bind:value={password} placeholder="Password">
    <button on:click={loginUser}>Login</button>
</div>

