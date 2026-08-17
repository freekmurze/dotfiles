## Effect Dependencies

### Never Suppress the Linter

```tsx
// BAD: Suppressing linter hides bugs
useEffect(() => {
  const id = setInterval(() => {
    setCount(count + increment);
  }, 1000);
  return () => clearInterval(id);
  // eslint-disable-next-line react-hooks/exhaustive-deps
}, []);

// GOOD: Fix the code, not the linter
useEffect(() => {
  const id = setInterval(() => {
    setCount(c => c + increment);
  }, 1000);
  return () => clearInterval(id);
}, [increment]);
```

### Use Updater Functions to Remove State Dependencies

```tsx
// BAD: messages in dependencies causes reconnection on every message
useEffect(() => {
  connection.on('message', (msg) => {
    setMessages([...messages, msg]);
  });
  // ...
}, [messages]); // Reconnects on every message!

// GOOD: Updater function removes dependency
useEffect(() => {
  connection.on('message', (msg) => {
    setMessages(msgs => [...msgs, msg]);
  });
  // ...
}, []); // No messages dependency needed
```

### Move Objects/Functions Inside Effects

```tsx
// BAD: Object created each render triggers Effect
function ChatRoom({ roomId }) {
  const options = { serverUrl, roomId }; // New object each render
  useEffect(() => {
    const connection = createConnection(options);
    connection.connect();
    return () => connection.disconnect();
  }, [options]); // Reconnects every render!
}

// GOOD: Create object inside Effect
function ChatRoom({ roomId }) {
  useEffect(() => {
    const options = { serverUrl, roomId };
    const connection = createConnection(options);
    connection.connect();
    return () => connection.disconnect();
  }, [roomId, serverUrl]); // Only reconnects when values change
}
```

### useEffectEvent for Non-Reactive Logic

```tsx
// BAD: theme change reconnects chat
function ChatRoom({ roomId, theme }) {
  useEffect(() => {
    const connection = createConnection(serverUrl, roomId);
    connection.on('connected', () => {
      showNotification('Connected!', theme);
    });
    connection.connect();
    return () => connection.disconnect();
  }, [roomId, theme]); // Reconnects on theme change!
}

// GOOD: useEffectEvent for non-reactive logic
function ChatRoom({ roomId, theme }) {
  const onConnected = useEffectEvent(() => {
    showNotification('Connected!', theme);
  });

  useEffect(() => {
    const connection = createConnection(serverUrl, roomId);
    connection.on('connected', () => {
      onConnected();
    });
    connection.connect();
    return () => connection.disconnect();
  }, [roomId]); // theme no longer causes reconnection
}
```

### Wrap Callback Props with useEffectEvent

```tsx
// BAD: Callback prop in dependencies
function ChatRoom({ roomId, onReceiveMessage }) {
  useEffect(() => {
    connection.on('message', onReceiveMessage);
    // ...
  }, [roomId, onReceiveMessage]); // Reconnects if parent re-renders
}

// GOOD: Wrap callback in useEffectEvent
function ChatRoom({ roomId, onReceiveMessage }) {
  const onMessage = useEffectEvent(onReceiveMessage);

  useEffect(() => {
    connection.on('message', onMessage);
    // ...
  }, [roomId]); // Stable dependency list
}
```

