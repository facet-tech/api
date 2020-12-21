function getDomPath(el) {
    // returns empty path for non valid element
    if (!isElement(el)) {
        return '';
    }
    var stack = [];
    while (el.parentNode != null) {
        var sibCount = 0;
        var sibIndex = 0;
        for (var i = 0; i < el.parentNode.childNodes.length; i++) {
            var sib = el.parentNode.childNodes[i];
            if (sib.nodeName == el.nodeName) {
                if (sib === el) {
                    sibIndex = sibCount;
                }
                sibCount++;
            }
        }
        if (el.hasAttribute('id') && el.id != '') {
            stack.unshift(el.nodeName.toLowerCase() + '#' + el.id);
        } else if (sibCount > 1) {
            stack.unshift(el.nodeName.toLowerCase() + ':eq(' + sibIndex + ')');
        } else {
            stack.unshift(el.nodeName.toLowerCase());
        }
        el = el.parentNode;
    }
    var res = stack.slice(1).join(' > '); // removes the html element
    return res.replace(/\s+/g, '');
}

function isElement(element) {
    return element instanceof Element || element instanceof HTMLDocument;
}

// const data = new Map([
//     {{.GO_ARRAY_REPLACE_ME}
// }
// ])

// used during debugging
const data = new Map([
    [window.location.pathname, new Set([
        'body>nav#ftco-navbar>div>div#ftco-nav', 'body>section:eq(0)>div>div:eq(1)',
        'body>section:eq(2)>div>div:eq(1)>div:eq(2)>a',
        'body>section:eq(3)>div>div:eq(1)>div>div>div:eq(0)>div>div:eq(3)>div>div>div:eq(0)',
        'body>section:eq(3)>div>div:eq(1)>div>div>div:eq(0)>div>div:eq(4)>div>div>div:eq(0)',
        'body>section:eq(3)>div>div:eq(1)>div>div>div:eq(0)>div>div:eq(2)>div>div>div:eq(1)>span',
        'body>section:eq(3)>div>div:eq(1)>div>div>div:eq(0)>div>div:eq(3)>div>div>div:eq(1)>p:eq(1)',
        'body>section:eq(3)>div>div:eq(1)>div>div>div:eq(0)>div>div:eq(4)>div>div>div:eq(1)>p:eq(0)',
        'body>section:eq(3)>div>div:eq(1)>div>div>div:eq(2)>button:eq(1)',
        'body>section:eq(3)>div>div:eq(1)>div>div>div:eq(0)>div>div:eq(3)>div>div>div:eq(1)>p:eq(0)',
        'body>section:eq(3)>div>div:eq(0)>div>h2', 'body>section:eq(3)>div>div:eq(0)>div>span',
        'body>footer>div>div:eq(0)>div:eq(0)', 'body>footer>div>div:eq(0)>div:eq(1)', 'body>footer>div>div:eq(0)>div:eq(3)',
        'body>footer>div>div:eq(0)>div:eq(2)', 'body>footer>div>div:eq(1)>div>p>i', 'body>section:eq(5)>div>div:eq(0)',
        'body>section#section-counter>div>div:eq(1)>div:eq(0)>div>div>span',
        'body>section:eq(3)>div>div:eq(1)>div>div>div:eq(0)>div>div:eq(2)>div>div>div:eq(1)>p:eq(1)',
        'body>section:eq(3)>div>div:eq(1)>div>div>div:eq(0)>div>div:eq(4)>div>div>div:eq(1)>p:eq(1)',
        'body>nav#ftco-navbar>div>div#ftco-nav', 'body>section:eq(0)>div>div:eq(1)', 'body>section:eq(0)>div>div:eq(0)',
        'body>nav#ftco-navbar>div>a', 'body>section:eq(1)>div>div>div', 'body>section:eq(2)>div>div:eq(0)',
        'body>section:eq(2)>div>div:eq(1)>div:eq(0)'])]
]);

var facetedNodes = new Set();
let nodesToRemove = data.get(window.location.pathname) || new Map();

/**
 * Computes whether the element's path is in the Set
 *
 * TODO: Consider using a different data structure (i.e: Trie) to improve performance
 */
const inHashMap = (element) => {
    const domPath = getDomPath(element);
    const inputSet = data.get(window.location.pathname);
    if (!inputSet) {
        return false;
    }
    let exists = false;
    inputSet.forEach(path => {
        // console.log('PATH', path, 'VS DOMPATH', domPath);
        if (domPath.includes(path)) {
            exists = true;
            return;
        }
    });
    return exists;
}

const callback = async function (mutationsList) {
    try {
        if (true) {
            for (let mutation of mutationsList) {
                // console.log('MUTATION', mutation)
                // TODO avoid iterating over subtrees that are not included
                if (mutation && mutation.target && mutation.target.children) {
                    const subPathContainedInMap = inHashMap(mutation.target);
                    if (subPathContainedInMap) {
                        continue;
                    }
                    domPathHide(mutation, mutation.target.children)
                }
            }
        }
    } catch (e) {
        console.log('[ERROR]', e);
    }
};

/**
 * Recursive function that iterates among DOM children
 *
 * @param {*} mutation
 * @param {*} mutationChildren
 */
const domPathHide = (mutation, mutationChildren) => {
    if (!mutationChildren) {
        return;
    }
    for (child of mutationChildren) {
        const childDomPath = getDomPath(child);
        if (nodesToRemove.has(childDomPath) && !facetedNodes.has(childDomPath)) {
            facetedNodes.add(childDomPath);
            child.style.display = "none";
            child.style.setProperty("display", "none", "important");
        }
        domPathHide(mutation, child.childNodes);
    }
}


const targetNode = document
const config = {subtree: true, childList: true, attributes: true};
const observer = new MutationObserver(callback);
observer.observe(targetNode, config);
