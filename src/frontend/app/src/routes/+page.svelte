<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>
<p>Welcome to the application. Please <a href="/login">log in</a> to continue.</p>
<script>
    import { goto } from '$app/navigation';
</script>

<!-- Using an anchor tag -->
<a href="/login">Go to Login</a>

<!-- Using a button with an event handler -->
<button on:click={() => goto('/login')}>Go to Login</button>
