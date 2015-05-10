/*document.title = ".edit";
document.body.innerHTML = "";*/

var files = {};
files.prefix = "/papill0n.org/shaders/";
files.makeName = function(name) {
  return `${name} - shaders!`;
}
/*files.current = "spec.txt";

files.builtin = {"spec.txt": `# Spec

- selecting a name from the list loads that file
- built-in files are accessible read-only
    (copy/paste if you want to change them)
- typing a new name and pressing enter (?) creates a new file
- typing an existing name and pressing enter loads that file

File store:

- \`.get(name: string) -> content or false if not existing\`
- \`create(name: string, content: string)\`, throws an error if
    a program with that name already exists (or if it is an
    internal file)
- \`save(name: string, content: string)\`, throws an error if
    there is no program with that name (?)
- \`remove(name: string)\`, for completeness

Some ideas:

- only store custom files in localStorage (obvious in retrospect)
- have a special temporary file that's saved very often (every
    keypress?), to prevent losing stuff.  reopen that file if
    it is present, delete/empty it on save.
- an 'unsaved changes' indicator, to make people save their stuff

Later:

- "forking" would be nice:
    - open existing code
    - say "fork" (or something more understandable)
    - create a new file with the content from the current one
    - edit away!
    - what about new names, though?
`,
             "default.frag": `void main() {
  gl_FragColor = gl_FragCoord;
}`,
             "includes/stdlib.frag": `float sdSphere(vec3 pos) {
  return length(pos) - 1.0;
}

float sdSphere(vec3 pos, float size) {
  return length(pos) - size;
}`,
             "includes/sphere-tracing.frag": `float trace(vec3 origin, vec3 direction) {
  // Uh, ... almost?
  return 0.0;
}`,
             "includes/ray-marching.frag": `// Warning: This is likely to be slower than sphere tracing!
float trace(vec3 origin, vec3 direction) {
  return 0.0;
}`};*/

files.exists = function(name) {
  return name in files.builtin || (files.prefix + name) in localStorage;
}

files.get = function(name) {
  if (name in files.builtin) {
    return {"name": name, "content": files.builtin[name], "readonly": true};
  } else {
    return {"name": name, "content": localStorage[files.prefix + name], "readonly": false};
  }
}

files.open = function(name) {
  files.current = name;
  document.title = files.makeName(name);
  files.currentFile = files.get(name);
  return files.currentFile;
}

files.create = function(name, content) {
  if (name.trim() == "") { throw new Error("name can't be empty"); }
  if (name in files.builtin) { throw new Error("can't use internal file names"); }
  if (files.exists(name)) { throw new Error("already exists"); }
  
  files.current = name;
  document.title = files.makeName(name);
  
  files.currentFile = { "name": name, "content": content, readonly: false };
  localStorage[files.prefix + name] = content;
  return files.currentFile;
}

files.save = function(name, content) {
  files.currentFile.content = content;
  localStorage[files.prefix + name] = content;
}

files.setupUI = function(editorEl) {
  function setFile(file, nameEl, contentEl) {
    nameEl.value = file.name;
    contentEl.value = file.content;
    contentEl.readOnly = file.readonly;
  }

  var nameEl = document.createElement("input");
  nameEl.id = "editor-name";
  nameEl.value = files.current;
  nameEl.maxlength = 20;
  nameEl.placeholder = "Name of the shader";
  nameEl.setAttribute('list', 'files');

  nameEl.onselect = function(ev) {
    setFile(files.open(nameEl.value), nameEl, editorEl);
    editorEl.focus();
  }

  nameEl.onkeyup = function(ev) {
    if (ev.keyCode == 13 && nameEl.value != files.current) { // Enter
      var file = nameEl.value;
      if (files.exists(file)) {
        setFile(files.open(file), nameEl, editorEl);
      } else {
        setFile(files.create(file, ""), nameEl, editorEl);
        var fileEl = document.createElement("option");
        fileEl.textContent = file;
        filesEl.appendChild(fileEl);
      }
      editorEl.focus();
    }
  }

  function addFileOption(name) {
    var fileEl = document.createElement("option");
    fileEl.textContent = name;
    filesEl.appendChild(fileEl);
  }

  var filesEl = document.createElement("datalist");
  filesEl.id = "files";
  Object.keys(files.builtin).forEach(addFileOption);

  Object.keys(localStorage).forEach(function(file) {
    if (file.startsWith(files.prefix)) {
      addFileOption(file.substr(files.prefix.length));
    }
  });

  var changeEl = document.createElement("span");
  changeEl.id = "editor-changed";
  
  editorEl.addEventListener("keydown", function(ev) {
    if (ev.ctrlKey && ev.keyCode == 83) { // Ctrl-s
      ev.preventDefault();

      changeEl.textContent = "";
      files.save(files.current, editorEl.value);
    }
  });

  editorEl.addEventListener("keyup", function(ev) {
    if (!ev.altKey && !ev.ctrlKey) {
      if (editorEl.value != files.currentFile.content) {
        changeEl.textContent = "(changed)";
      } else {
        changeEl.textContent = "";
      }
    }
  });

  setFile(files.open(files.current), nameEl, editorEl);

  var containerEl = document.createElement("div");
  containerEl.appendChild(nameEl);
  containerEl.appendChild(filesEl);
  containerEl.appendChild(changeEl);
  return containerEl;
}

/*var editor = {};
editor.el = document.createElement("textarea");
editor.el.style = "width: 70ex; height: 40em";
editor.nameUIEl = files.setupUI(editor.el);

document.body.appendChild(editor.nameUIEl);
document.body.appendChild(editor.el);*/