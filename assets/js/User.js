class User {
	async createUser(superSecretPassword) {
		const resp = await fetch(window.createUser, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ password: superSecretPassword }),
			credentials: "include",
		});
		const data = await resp.json();
		console.log(data,' data');
		if (data.isSuccess) {
			alert('Юзер был успешно создан');
			const cookie = new Cookie();
			cookie.set('user_id', data.newUserId, 360);
			window.location.reload();
		} else {
			alert('Не удалось создать юзера');
		}
	}
}
