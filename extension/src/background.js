const MINUTE_IN_MS = 1000 * 60;

(async function updateMutedTwitterAccounts() {
  try {
    const response = await fetch('http://localhost:8080/muted-accounts');
    const mutedTwitterAccounts = await response.json();
    chrome.storage.local.set({ mutedTwitterAccounts });

    const [tab] = await chrome.tabs.query({ active: true, lastFocusedWindow: true, url: 'https://x.com/*' });
    if (!tab) {
      return;
    }

    chrome.scripting.executeScript({
      target: { tabId: tab.id },
      function: (mutedTwitterAccounts) => {
        window.postMessage({ type: "MUTEX_EXTENSION", message: mutedTwitterAccounts });
      },
      args: [mutedTwitterAccounts]
    });
  } catch (error) {
    console.error(error);
  } finally {
    setTimeout(updateMutedTwitterAccounts, 15 * MINUTE_IN_MS);
  }
})()

chrome.runtime.onMessageExternal.addListener((message, sender, sendResponse) => {
  if (message.type === 'getMutedTwitterAccounts') {
    chrome.storage.local.get('mutedTwitterAccounts', ({ mutedTwitterAccounts }) => {
      sendResponse(mutedTwitterAccounts);
    });
    return true;
  } else if (message.type === 'muteTwitterAccount') {
    fetch('http://localhost:8080/mute', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ user_handle: message.userHandle })
    }).then(response => {
      sendResponse(response.ok);
    }).catch((err) => {
      console.error(err);
      sendResponse(false) 
    });
  } else if (message.type === 'addCachedMutedTwitterAccount') {
    chrome.storage.local.get('mutedTwitterAccounts', ({ mutedTwitterAccounts }) => {
      mutedTwitterAccounts.push(message.userHandle);
      chrome.storage.local.set({ mutedTwitterAccounts });
    });
  }
});

