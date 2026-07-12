#!/usr/bin/env python3
"""sync-agents.py — 从 .agents/agents/ 真源同步到四端。

用法: python3 .agents/skills/sync-agents/scripts/sync-agents.py
从仓库根执行。
"""
import os, re, sys

ROOT = os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))))
SRC = os.path.join(ROOT, '.agents/agents')
TARGETS = {
    '.claude/agents': {'tools': lambda s: ' '.join(w.capitalize() for w in s.split()) if s else s},
    '.opencode/agents': {'tools': lambda s: s},
    '.pi/agents': {'tools': lambda s: s},
}
CODEX_DIR = os.path.join(ROOT, '.codex/agents')


def parse(path):
    with open(path) as f:
        content = f.read()
    m = re.match(r'^---\s*\n(.*?)\n---\s*\n(.*)', content, re.DOTALL)
    if not m:
        return None
    fm = {}
    for line in m.group(1).splitlines():
        if ':' in line:
            k, _, v = line.partition(':')
            fm[k.strip()] = v.strip()
    return fm, m.group(2)


def main():
    for root, dirs, files in os.walk(SRC):
        for fname in sorted(files):
            if not fname.endswith('.md'):
                continue
            src_path = os.path.join(root, fname)
            res = parse(src_path)
            if not res:
                continue
            fm, body = res
            name = fname[:-3]
            desc = fm.get('description', '')
            tools = fm.get('tools', 'read, grep, glob')
            model = fm.get('model', 'inherit')
            skills = fm.get('skills', '')

            # MD 端
            for tdir, transforms in TARGETS.items():
                lines = [
                    '---',
                    f'name: {name}',
                    f'description: {desc}',
                    f'tools: {transforms["tools"](tools)}',
                    f'model: {model}',
                ]
                if skills:
                    lines.append(f'skills: {skills}')
                lines.append('---')
                lines.append('')
                lines.append(body)
                out_dir = os.path.join(ROOT, tdir)
                os.makedirs(out_dir, exist_ok=True)
                with open(os.path.join(out_dir, fname), 'w') as f:
                    f.write('\n'.join(lines))

            # Codex TOML
            os.makedirs(CODEX_DIR, exist_ok=True)
            esc_body = body.replace('"', '\\"')
            with open(os.path.join(CODEX_DIR, f'{name}.toml'), 'w') as f:
                f.write(f'name = "{name}"\n')
                f.write(f'description = """{desc}"""\n')
                f.write(f'developer_instructions = """\n{body}\n"""\n')

    count = len(os.listdir(SRC))
    claude_count = len(os.listdir(os.path.join(ROOT, '.claude/agents')))
    opencode_count = len(os.listdir(os.path.join(ROOT, '.opencode/agents')))
    pi_count = len(os.listdir(os.path.join(ROOT, '.pi/agents')))
    codex_count = len(os.listdir(CODEX_DIR))
    print(f"Synced {count} agents from .agents/agents/")
    print(f"  .claude/agents:   {claude_count}")
    print(f"  .opencode/agents: {opencode_count}")
    print(f"  .pi/agents:       {pi_count}")
    print(f"  .codex/agents:    {codex_count}")


if __name__ == '__main__':
    main()
