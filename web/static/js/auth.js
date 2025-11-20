const EMAIL_KEY = 'vado_user_email'

export async function initAuth() {
    try {
        const cachedEmail = sessionStorage.getItem(EMAIL_KEY);
        if (cachedEmail) {
            return { email: cachedEmail };
        }

        const me = await fetch("/me").then(r => {
            if (r.ok) {
                return r.json();
            }
        })

        if (me && me.email) {
            sessionStorage.setItem(EMAIL_KEY, me.email);
        }

        return me
    } catch (e) {
        console.log("Error", e)
    }
}

export async function logout() {
    sessionStorage.removeItem(EMAIL_KEY)
}