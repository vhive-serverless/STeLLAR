export default {
	async fetch(request) {
		var incrLimit = 0

		if (request.url.includes("incrementLimit")) {
			incrLimit = parseInt(request.url.split("=")[1])
		}

		simulateWork(incrLimit)

		const resData = {
			"RequestID": "cloudflare-does-not-specify",
			"TimestampChain": [Date.now().toString()],
		};

		const body = JSON.stringify(resData, null, 2);

		return new Response(body, {
			headers: {
				"content-type": "application/json;charset=UTF-8",
			},
		});
	},
};

function simulateWork(incr) {
	let i = 0;
	while (i < incr) {
		i++;
	}
}

