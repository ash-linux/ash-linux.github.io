const vscode = require('vscode');
const http = require('http');

function activate(context) {
    let disposable = vscode.commands.registerCommand('ash.search', async function () {
        const query = await vscode.window.showInputBox({
            placeHolder: 'Semantic search query...'
        });
        
        if (!query) return;
        
        const reqData = JSON.stringify({ query: query, limit: 10 });
        
        const options = {
            hostname: 'localhost',
            port: 8080,
            path: '/api/search',
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Content-Length': reqData.length
            }
        };
        
        const req = http.request(options, (res) => {
            let data = '';
            res.on('data', (chunk) => data += chunk);
            res.on('end', async () => {
                try {
                    const parsed = JSON.parse(data);
                    const items = parsed.results.map(r => ({
                        label: r.path.split('/').pop(),
                        description: r.path,
                        detail: `Score: ${r.score}`
                    }));
                    
                    const selected = await vscode.window.showQuickPick(items);
                    if (selected) {
                        const doc = await vscode.workspace.openTextDocument(selected.description);
                        vscode.window.showTextDocument(doc);
                    }
                } catch (e) {
                    vscode.window.showErrorMessage('Ash Search failed to parse response.');
                }
            });
        });
        
        req.on('error', (e) => {
            vscode.window.showErrorMessage('Ash API not reachable. Is ash-api running?');
        });
        
        req.write(reqData);
        req.end();
    });

    context.subscriptions.push(disposable);
}

function deactivate() {}

module.exports = {
    activate,
    deactivate
}
