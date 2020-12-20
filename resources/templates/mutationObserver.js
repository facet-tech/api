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

const data = new Map([
    {{.GO_ARRAY_REPLACE_ME}}
])

var facetedNodes = new Set();
let nodesToRemove = data.get(window.location.pathname) || new Map();

const callback = async function (mutationsList) {
    try {
        if (data.has(window.location.pathname)) {
            for (let mutation of mutationsList) {
                // TODO avoid iterating over subtrees that are not included
                if (mutation && mutation.target && mutation.target.children) {
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
            child.style.display = "none"
            child.style.setProperty("display", "none", "important");
        }
        domPathHide(mutation, child.childNodes);
    }
}


const targetNode = document
const config = {subtree: true, childList: true, attributes: true};
const observer = new MutationObserver(callback);
observer.observe(targetNode, config);
