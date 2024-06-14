import { render, html } from 'lit-html';

const extensionId = 'lnobhlnamiibdgjdaojoolicmodaohlf';

let mutedTwitterAccounts = new Set();

/**
 * 
 * @param {MutationRecord} mutation 
 */
function handleDeleteMutation(mutation) {
  mutation.addedNodes.forEach(node => {
    if (node.nodeType !== node.ELEMENT_NODE) {
      return;
    }
  
    const tweet = node.querySelector('[data-testid="tweet"]');
    if (!tweet) {
      return;
    }
  
    const anchor = tweet.querySelector('[data-testid="User-Name"] a');
    if (!anchor) {
      return;
    }

    const userHandle = anchor.href.split('/').at(-1);

    if (mutedTwitterAccounts.has(userHandle)) {
      tweet.closest('[data-testid="cellInnerDiv"]')?.remove();
    }
  });
}

/**
 * 
 * @param {MutationRecord} mutation 
 */
function handleAddMuteButtonMutation(mutation) {
  mutation.addedNodes.forEach(node => {
    if (node.nodeType !== node.ELEMENT_NODE) {
      return;
    }
  
    /**
     * @type {HTMLDivElement}
     */
    const dropdown = node.querySelector('[data-testid="Dropdown"]');
    if (!dropdown) {
      return;
    }

    const result = /\@\w+/.exec(dropdown.innerText);

    if (!result) {
      return;
    }

    const userHandle = result[0].replace('@', '');

    render(muteButtonTemplate(userHandle), dropdown, { renderBefore: dropdown.firstChild });
  });
}

chrome.runtime.sendMessage(extensionId, { type: 'getMutedTwitterAccounts' }, mutedAccounts => {
  mutedTwitterAccounts = new Set(mutedAccounts);
});

const deleteTweetObserver = new MutationObserver(mutations => mutations.forEach(handleDeleteMutation));
const muteTweetObserver = new MutationObserver(mutations => mutations.forEach(handleAddMuteButtonMutation));

navigation.addEventListener('navigate', event => {
  const isTweetRgx = /https:\/\/x.com\/\w+\/status\/\d+/;
  deleteTweetObserver.disconnect();
  muteTweetObserver.disconnect();
  if (isTweetRgx.test(event.destination.url)) {
    const twitterApp = document.getElementById('react-root');
    if (!twitterApp) {
      return;
    }

    deleteTweetObserver.observe(twitterApp, { subtree: true, childList: true });
    muteTweetObserver.observe(twitterApp, { subtree: true, childList: true });
  }
});

window.addEventListener("message", (event) => {
  if (event.source !== window) {
    return;
  }

  if (event.data.type === "MUTEX_EXTENSION_ID") {
    console.log('MUTEX_EXTENSION_ID call', event.data.message);
  }

  if (event.data?.type === "MUTEX_EXTENSION") {
    console.log('MUTEX_EXTENSION call');
    mutedTwitterAccounts = new Set(event.data.message);
  }
});

async function muteBot(userHandle) {
  mutedTwitterAccounts.add(userHandle);
  chrome.runtime.sendMessage(extensionId, { type: 'muteTwitterAccount', userHandle }, success => {
    if (!success) {
      mutedTwitterAccounts.delete(userHandle);
      return alert(`An error occurred while trying to mute ${userHandle}`);
    }

    document
      .querySelectorAll(`[data-testid="tweet"] a[href="/${userHandle}"]`)
      ?.forEach(anchor => {
        anchor.closest('[data-testid="cellInnerDiv"]')?.remove();
      });
    
    document
      .querySelector('.css-175oi2r.r-1p0dtai.r-1d2f490.r-1xcajam.r-zchlnj.r-ipm5af')
      ?.click();
    
    chrome.runtime.sendMessage(extensionId, { type: 'addCachedMutedTwitterAccount', userHandle });
  });
}

function muteButtonTemplate(userHandle) {
  return html`
    <div .onclick="${() => muteBot(userHandle)}" role="menuitem" tabindex="0" class="css-175oi2r r-1loqt21 r-18u37iz r-1mmae3n r-3pj75a r-13qz1uu r-o7ynqc r-6416eg r-1ny4l3l">
      <div class="css-175oi2r r-1777fci r-faml9v">
        <svg viewBox="0 0 24 24" aria-hidden="true" class="r-4qtqp9 r-yyyyoo r-1xvli5t r-dnmrzs r-bnwqim r-lrvibr r-m6rgpd r-1nao33i r-1q142lx"><g><path d="M18 6.59V1.2L8.71 7H5.5C4.12 7 3 8.12 3 9.5v5C3 15.88 4.12 17 5.5 17h2.09l-2.3 2.29 1.42 1.42 15.5-15.5-1.42-1.42L18 6.59zm-8 8V8.55l6-3.75v3.79l-6 6zM5 9.5c0-.28.22-.5.5-.5H8v6H5.5c-.28 0-.5-.22-.5-.5v-5zm6.5 9.24l1.45-1.45L16 19.2V14l2 .02v8.78l-6.5-4.06z"></path></g></svg>
      </div>
      <div class="css-175oi2r r-16y2uox r-1wbh5a2">
        <div dir="ltr" class="r-bcqeeo r-1ttztb7 r-qvutc0 r-37j5jr r-a023e6 r-rjixqe r-b88u0q" style="text-overflow: unset; color: rgb(231, 233, 234);">
          <span class="css-1jxf684 r-bcqeeo r-1ttztb7 r-qvutc0 r-poiln3" style="text-overflow: unset;">Block with MuteX @${userHandle}</span>
        </div>
      </div>
    </div>`
}
