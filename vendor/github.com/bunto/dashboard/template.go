package dashboard

import (
	"html/template"
)

type templateInfo struct {
	Projects []*Project
}

var (
	indexTmpl = template.Must(template.New("index.html").Parse(`
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
  <meta content="origin-when-cross-origin" name="referrer" />
  <link crossorigin="anonymous" href="https://assets-cdn.github.com/assets/frameworks-98550932b9f11a849da143d2dbc9dfaa977a17656514d323ae9ce0d6fa688b60.css" integrity="sha256-mFUJMrnxGoSdoUPS28nfqpd6F2VlFNMjrpzg1vpoi2A=" media="all" rel="stylesheet" />
  <link crossorigin="anonymous" href="https://assets-cdn.github.com/assets/github-a17b8c9d020ded73daa7ee1a3844b4512f12076634d9249861ddf86dc33da66e.css" integrity="sha256-oXuMnQIN7XPap+4aOES0US8SB2Y02SSYYd34bcM9pm4=" media="all" rel="stylesheet" />
  <title>Dashboard</title>
  <style type="text/css">
  .markdown-body {
      width: 95%;
      margin: 0 auto;
  }
  </style>
  <script type="application/javascript">
  function reqListener () {
    console.log(this.responseText);

    if (this.responseText === null || this.responseText === "") {
        console.log("nada");
        return;
    }

    var info = JSON.parse(this.responseText);
    var tr = document.getElementById(info.name);

    // Name
    var nameTD = document.createElement("td");
    var nameA = document.createElement("a");
    nameA.href = "https://github.com/" + info.nwo;
    nameA.title = info.name + " on GitHub";
    nameA.innerText = info.name;
    nameTD.appendChild(nameA);
    tr.appendChild(nameTD);

    // Gem
    var rubygemsTD = document.createElement("td");
    if (info.gem) {
        var rubygemsA = document.createElement("a");
        rubygemsA.href = "https://rubygems.org/gems/" + info.gem.name;
        rubygemsA.title = info.gem.name + " on RubyGems.org";
        rubygemsA.innerText = "v" + info.gem.version;
        rubygemsTD.appendChild(rubygemsA);
    } else {
        rubygemsTD.innerText = "no info";
    }
    tr.appendChild(rubygemsTD);

    // Travis
    var travisTD = document.createElement("td");
    if (info.travis) {
        var travisA = document.createElement("a");
        travisA.href = "https://travis-ci.org/" + info.travis.nwo + "/builds/" + info.travis.branch.id;
        travisA.title = info.travis.nwo + " on TravisCI";
        travisA.innerText = info.travis.branch.state;
        travisTD.appendChild(travisA);
    } else {
        travisTD.innerText = "no info";
    }
    tr.appendChild(travisTD);

    // Downloads
    var downloadsTD = document.createElement("td");
    if (info.gem && info.gem.downloads) {
        downloadsTD.innerText = info.gem.downloads;
    } else {
        downloadsTD.innerText = "no info";
    }
    tr.appendChild(downloadsTD);

    if (info.github === undefined || info.github === null) {
        for (i = 0; i < 4; i++) {
            var emptyTD = document.createElement("td");
            emptyTD.innerText = "no info";
            tr.appendChild(emptyTD);
        }
        return;
    }

    // Commits
    var commitsTD = document.createElement("td");
    commitsTD.innerText = info.github.commits_this_week;
    tr.appendChild(commitsTD);

    // Pull Requests
    var pullrequestsTD = document.createElement("td");
    if (info.github.open_prs > 0) {
        var pullrequestsA = document.createElement("a");
        pullrequestsA.href = "https://github.com/" + info.nwo + "/pulls";
        pullrequestsA.title = info.name + " pull requests on GitHub";
        pullrequestsA.innerText = info.github.open_prs;
        pullrequestsTD.appendChild(pullrequestsA);
    } else if (info.github.open_prs < 0) {
        pullrequestsTD.innerText = "no info";
    } else {
        pullrequestsTD.innerText = info.github.open_prs;
    }
    tr.appendChild(pullrequestsTD);

    // Issues
    var issuesTD = document.createElement("td");
    if (info.github.open_issues > 0) {
        var issuesA = document.createElement("a");
        issuesA.href = "https://github.com/" + info.nwo + "/issues";
        issuesA.title = info.name + " issues on GitHub";
        issuesA.innerText = info.github.open_issues;
        issuesTD.appendChild(issuesA);
    } else if (info.github.open_issues < 0) {
        issuesTD.innerText = "no info";
    } else {
        issuesTD.innerText = info.github.open_issues;
    }
    tr.appendChild(issuesTD);

    // Unreleased commits
    var unreleasedcommitsTD = document.createElement("td");
    if (info.github.commits_since_latest_release > 0) {
        var unreleasedcommitsA = document.createElement("a");
        unreleasedcommitsA.href = "https://github.com/" + info.nwo + "/compare/" + info.github.latest_release_tag + "...master";
        unreleasedcommitsA.title = info.name + " commits since latest release on GitHub";
        unreleasedcommitsA.innerText = info.github.commits_since_latest_release;
        unreleasedcommitsTD.appendChild(unreleasedcommitsA);
    } else if (info.github.commits_since_latest_release < 0) {
        unreleasedcommitsTD.innerText = "no info";
    } else {
        unreleasedcommitsTD.innerText = info.github.commits_since_latest_release;
    }
    tr.appendChild(unreleasedcommitsTD);
  }

  {{range .Projects}}
  var oReq = new XMLHttpRequest();
  oReq.addEventListener("load", reqListener);
  oReq.open("GET", "/show.json?name={{urlquery .Name}}");
  oReq.send();
  {{end}}
  </script>
</head>
<body>
<div class="markdown-body">

<table>
  <caption>Bunto At-a-Glance Dashboard</caption>
  <thead>
    <tr>
      <th>Repo</th>
      <th>Gem</th>
      <th>Travis</th>
      <th>Downloads</th>
      <th>Commits</th>
      <th>Pull Requests</th>
      <th>Issues</th>
      <th>Unreleased commits</th>
    </tr>
  </thead>
  <tbody>
    {{range .Projects}}
    <tr id="{{.Name}}"></tr>
    {{end}}
  </tbody>
</table>

<div>
	<strong>Commits are as of this week. Issues and pull requests are total open.</strong>
	<a href="https://github.com/bunto/dashboard">Source Code</a>.
</div>

<p>
	Look wrong? <form action="/reset.json" method="post"><input type="Submit" value="Reset the cache."></form>
</p>

</div>
</body>
</html>
`))
)
