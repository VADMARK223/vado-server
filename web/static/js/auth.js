export async function initAuth() {
    try {
        const cachedEmail = sessionStorage.getItem('userEmail');
        if (cachedEmail) {
            return { email: cachedEmail };
        }

        const me = await fetch("/me").then(r => {
            if (r.ok) {
                return r.json();
            }
        })

        if (me && me.email) {
            sessionStorage.setItem('userEmail', me.email);
        }

        return me
    } catch (e) {
        console.log("Error", e)
    }
}