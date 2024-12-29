/**
 * Handle 3 different cases
 * 
 * 1. Already focusing a chrome tab, then change to another tab (chrome.tabs.onActivated)
 * 2. Already focusing a chrome tab, then change the current tab to another url (chrome.tabs.onUpdated)
 * 3. Not focusing a chrome tab (e.g. another app), then focus a chrome tab (chrome.windows.onFocusChanged)
 */


const eligibleForListener = (tab: chrome.tabs.Tab) => tab.incognito === false && tab.title !== 'New Tab' && !tab.url?.startsWith('chrome://')

const getTabInfo = (tab: chrome.tabs.Tab) => ({
  'app_name': 'Google Chrome',
  'window_title': tab.title,
  'url': tab.url,
  'bundle_identifier': 'com.google.Chrome'
});

const handleTabChange = (tab: chrome.tabs.Tab) => {
  if (!eligibleForListener(tab)) {
    return;
  }

  const tabInfo = getTabInfo(tab);
  fetch('http://localhost:6969/event', {
    method: 'POST', // Use POST method
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(tabInfo)
  })
  .then(response => response.status)
  .then(status => {
    if (status !== 200) {
      console.error('Error making request');
    }
  })
  .catch(error => {
    console.error('Error making POST request:', error);
  });
}

chrome.tabs.onActivated.addListener((activeInfo) => {
  chrome.tabs.get(activeInfo.tabId, (tab) => {
    handleTabChange(tab);
  });
});

chrome.tabs.onUpdated.addListener((tabId, changeInfo, tab) => {
  // Only proceed if the status is 'complete' (finished loading the page)
  // Or the title changed (for YouTube)

  // Most websites: change title, and all metadata -> complete
  // For YouTube: update some metadata -> complete -> change title

  // Another edge case: YouTube front page -> click video -> complete -> change title to previous video -> change title to current video
  // because why the f**k not?
  if (changeInfo.status === 'complete' || changeInfo.title !== undefined) {
    chrome.tabs.query({ active: true, currentWindow: true }, function(tabs) {
      if (tabs.length > 0) {
        // Get the first (and only) active tab
        const activeTab = tabs[0];
        if (activeTab.id !== tabId) {
          return;
        }

        // this is not the last event; wait for title change
        if (tab.url?.includes('youtube.com') && changeInfo.status === 'complete') {
          return;
        }

        if (!tab.url?.includes('youtube.com') && changeInfo.status !== 'complete') {
          return;
        }

        handleTabChange(tab);
      } else {
        console.log("No active tab found.");
      }
    });
  }
});

chrome.windows.onFocusChanged.addListener((windowId) => {
  if (windowId !== chrome.windows.WINDOW_ID_NONE) {
    // console.log('Chrome window is focused');
    // console.log(windowId);
    chrome.tabs.query({ windowId: windowId, active: true }, function(tabs) {
      if (tabs.length > 0) {
        let activeTab = tabs[0];
        handleTabChange(activeTab);
      } else {
        console.log('No active tab in this window.');
      }
    });
  }
});
