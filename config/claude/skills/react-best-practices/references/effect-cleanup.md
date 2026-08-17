## Effect Cleanup

### Always Clean Up Subscriptions

```tsx
useEffect(() => {
  const connection = createConnection(serverUrl, roomId);
  connection.connect();
  return () => connection.disconnect(); // REQUIRED
}, [roomId]);

useEffect(() => {
  function handleScroll(e) {
    console.log(window.scrollY);
  }
  window.addEventListener('scroll', handleScroll);
  return () => window.removeEventListener('scroll', handleScroll); // REQUIRED
}, []);
```

### Data Fetching with Ignore Flag

```tsx
useEffect(() => {
  let ignore = false;

  async function fetchData() {
    const result = await fetchTodos(userId);
    if (!ignore) {
      setTodos(result);
    }
  }

  fetchData();

  return () => {
    ignore = true; // Prevents stale data from old requests
  };
}, [userId]);
```

### Development Double-Fire Is Intentional

React remounts components in development to verify cleanup works. If you see effects firing twice, don't try to prevent it with refs:

```tsx
// BAD: Hiding the symptom
const didInit = useRef(false);
useEffect(() => {
  if (didInit.current) return;
  didInit.current = true;
  // ...
}, []);

// GOOD: Fix the cleanup
useEffect(() => {
  const connection = createConnection();
  connection.connect();
  return () => connection.disconnect(); // Proper cleanup
}, []);
```

