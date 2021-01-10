const Controller = {
  fetchFilters: (ev) => {
    ev.preventDefault();
    const response = fetch(`/filters`).then((response) => {
      response.json().then((results) => {
        filter = document.getElementById('filter');
        clearOptions(filter);
        filter.options.add(new Option("ALL", "ALL"));
        $.each(results, function (index, item) {
          filter.options.add(new Option(item, item));
        });
      });
    });
  },

  search: (ev) => {
    ev.preventDefault();
    const form = document.getElementById("form");
    const data = Object.fromEntries(new FormData(form));
    const response = fetch(`/search?q=${data.query}&filter=${data.filter}`).then((response) => {
      response.json().then((results) => {
        Controller.renderResult(results, data.query);
      });
    });
  },

  renderResult: (results, query) => {
    $('#pagination-container').pagination({
      dataSource: results,
      pageSize: 5,
      callback: function (data, pagination) {
        var container = $('#pagination-container');
        var dataHtml = '<ul>';

        $.each(data, function (index, item) {
          result = highlightQuery(item, query);
          result = insertEllipsis(result);
          dataHtml += '<li><pre>' + result + '</pre></li>';
        });

        dataHtml += '</ul>';

        container.prev().html(dataHtml);
      }
    })
  },
}

function highlightQuery(result, query) {
  var regex = new RegExp("(" + query + ")", "gi");
  return result.replace(regex, '<b>$1</b>')
}

function insertEllipsis(result) {
  result = result.trim();
  if (result.endsWith('?') || result.endsWith('.') || result.endsWith('!')) {
    return result;
  }

  return result + " <b>...</b>";
}

function clearOptions(selectElement) {
  var i, L = selectElement.options.length - 1;
  for (i = L; i >= 0; i--) {
    selectElement.remove(i);
  }
}

window.onload = Controller.fetchFilters

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
