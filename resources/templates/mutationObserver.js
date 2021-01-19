/**
 * Returns the DOM path of the element
 *
 * @param el
 * @returns {string}
 *
 */
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

/**
 * Returns whether a given element in an HTML element or not
 *
 * @param element
 * @returns {boolean}
 */
function isElement(element) {
    return element instanceof Element || element instanceof HTMLDocument;
}

const globalFacetKey = 'GLOBAL-FACET-DECLARATION';
const data = {{.GO_ARRAY_REPLACE_ME}};

let nodesToRemove = (data[window.location.pathname] || []).concat(data[globalFacetKey] || []) || [];
console.log('nodesToRemove', nodesToRemove);

let pathsAlreadyRemoved = [];

/**
 * Computes whether the element's path is in the Set
 *
 * @param {*} element
 */
const inHashMap = (element) => {
    const domPath = getDomPath(element);
    let exists = false;
    nodesToRemove.forEach(nodeElement => {
        if (nodeElement.path.includes(domPath)) {
            exists = true;
            return;
        }
    });
    return exists;
}

const callback = async function (mutationsList) {
    try {
        for (let mutation of mutationsList) {
            if (mutation && mutation.target && mutation.target.children) {
                const subPathContainedInMap = inHashMap(mutation.target);
                if (!subPathContainedInMap) {
                    continue;
                }
                domPathHide(mutation, mutation.target.children)
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
    for (const child of mutationChildren) {
        const childDomPath = getDomPath(child);
        const wantedElement = nodesToRemove.filter(e => e.path === childDomPath);
        if (!wantedElement || !wantedElement[0] || pathsAlreadyRemoved.includes(wantedElement[0].path)) {
            domPathHide(mutation, child.childNodes);
            continue;
        }
        if (wantedElement[0].domRemove) {
            child.remove()
            pathsAlreadyRemoved.push(childDomPath);
            continue;
        } else {
            child.style.display = "none";
            child.style.setProperty("display", "none", "important");
        }
        domPathHide(mutation, child.childNodes);
    }
}

const targetNode = document
const config = {subtree: true, childList: true, attributes: true};

/*
 * disableMutationObserver can be passed through the facet-extension to override this behavior
 */
if (typeof window.disableMutationObserverScript === 'undefined' || window.disableMutationObserverScript === undefined) {
    const observer = new MutationObserver(callback);
    observer.observe(targetNode, config);
} else {
    console.log('Facet extension is enabled. Blocking script execution.');
}
