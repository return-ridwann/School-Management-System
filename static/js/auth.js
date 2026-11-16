// check user auith in their role
async function checkAuth(expectedRole) {
    try {
        const res = await fetch('/api/me');
        if (!res.ok) {
            window.location.href = '/login.html';
            return null;
        }

        const data = await res.json();
        const user = {
            id: data.user_id || data.id,
            full_name: data.full_name,
            role: data.role
        };

        // if role dosen't match redirect to dashboard
        if (expectedRole && user.role !== expectedRole) {
            if (user.role === 'admin') {
                window.location.href = '/admin.html';
            } else if (user.role === 'teacher') {
                window.location.href = '/teacher.html';
            } else if (user.role === 'student') {
                window.location.href = '/student.html';
            } else {
                window.location.href = '/login.html';
            }
            return null;
        }

        return user;
    } catch (err) {
        console.error('Auth verification failed:', err);
        window.location.href = '/login.html';
        return null;
    }
}

// logout user and redirect to login page
async function logout() {
    try {
        await fetch('/api/logout', {
            method: 'POST'
        });
    } catch (err) {
        console.error('Logout error:', err);
    }
    window.location.href = '/login.html';
}

