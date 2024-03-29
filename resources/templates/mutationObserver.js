(function () {
    function getDomPath(el) {
        if (!el || !isElement(el)) {
            return '';
        }
        var stack = [];
        var isShadow = false;
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
            var nodeName = el.nodeName.toLowerCase();
            if (isShadow) {
                nodeName += "::shadow";
                isShadow = false;
            }
            if (sibCount > 1) {
                if (sibIndex === 0) {
                    stack.unshift(nodeName);
                } else {
                    stack.unshift(nodeName + ':nth-of-type(' + (sibIndex + 1) + ')');
                }
            } else {
                stack.unshift(nodeName);
            }
            el = el.parentNode;
            if (el.nodeType === 11) {
                isShadow = true;
                el = el.host;
            }
        }
        var res = stack.slice(1).join(' > ');
        return res.replace(/ /g, "");
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

    /**
     * Computes whether the element's path is in the nodesToRemove array
     *
     * @param {*} element
     */
    const inHashMap = (element) => {
        const domPath = getDomPath(element);
        let exists = false;
        nodesToRemove.forEach(path => {
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
            if (nodesToRemove.includes(childDomPath) && child.style) {
                child.style.display = "none";
                child.style.setProperty("display", "none", "important");
            }
            domPathHide(mutation, child.childNodes);
        }
    }

    function getFacetExtensionCookie(key) {
        const value = `; ${document.cookie}`;
        const parts = value.split(`; ${key}=`);
        if (parts.length === 2) {
            return parts.pop().split(';').shift();
        }
    }

    const facetExtensionDisableMo = getFacetExtensionCookie("FACET_EXTENSION_DISABLE_MO") === 'true';
    console.log('facetExtensionDisableMo',facetExtensionDisableMo);

    if(!facetExtensionDisableMo) {
        const config = {subtree: true, childList: true, attributes: true};
        const observer = new MutationObserver(callback);
        observer.observe(document, config);
    } else {
        console.log('[Facet][API][mutationObserver.js] Script Disabled');
    }
})();

