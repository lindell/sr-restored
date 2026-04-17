import { createServer } from 'node:http';
import { readFileSync, existsSync, statSync } from 'node:fs';
import { join, extname } from 'node:path';

const dir = 'build';
const mimeTypes = {
	'.html': 'text/html',
	'.js': 'text/javascript',
	'.css': 'text/css',
	'.json': 'application/json',
	'.png': 'image/png',
	'.svg': 'image/svg+xml',
	'.webp': 'image/webp',
	'.woff2': 'font/woff2',
	'.woff': 'font/woff'
};

createServer((req, res) => {
	let filePath = join(dir, new URL(req.url, 'http://localhost').pathname);
	if (existsSync(filePath) && statSync(filePath).isDirectory()) {
		filePath = join(filePath, 'index.html');
	}
	try {
		const data = readFileSync(filePath);
		const contentType = mimeTypes[extname(filePath)] || 'application/octet-stream';
		res.writeHead(200, { 'Content-Type': contentType });
		res.end(data);
	} catch {
		res.writeHead(404);
		res.end('Not found');
	}
}).listen(4173, () => {
	console.log('Static test server running on http://localhost:4173');
});
