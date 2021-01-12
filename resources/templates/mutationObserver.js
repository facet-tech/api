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

const data = new Map([
    {{.GO_ARRAY_REPLACE_ME}}
]);

/**
 * IIME for transforming domainpath-specific data into global facets
 */
(() => {
    let result = [];
    for (const [_, value] of data.entries()) {
        const arr = Array.from(value);
        result.push(arr);
    }
    transformedData = [].concat.apply([], result);
})();

let nodesToRemove = data.get(window.location.pathname) || new Map();

/**
 * Computes whether the element's path is in the Set
 *
 * @param {*} element
 */
const inHashMap = (element) => {
    const domPath = getDomPath(element);
    transformedData.forEach(path => {
        if (path.includes(domPath)) {
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
        if (transformedData.includes(childDomPath) && child.style) {
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
