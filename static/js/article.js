function init_editor(editor_id){
    CKEDITOR.replace( editor_id, {
        height: 250,
        extraPlugins: 'colorbutton',
        colorButton_enableAutomatic: false,
        codeSnippet_theme: 'zenburn',
        //toolbar : [ ['Bold','Italic','Underline','Strike','JustifyLeft','JustifyCenter','JustifyRight','Link'], ['Image','Table','HorizontalRule','SpecialChar'], ['TextColor','BGColor','RemoveFormat','Font','FontSize','Source'] ]
    } );

    window.CKEDITOR = CKEDITOR;

    CKEDITOR.instances[editor_id].on('change', function() {
        $('#'+editor_id).val(CKEDITOR.instances[editor_id].getData());
    });

}

function getParameter(name) {
    var reg = new RegExp("(^|&)"+ name +"=([^&]*)(&|$)");
    var r = window.location.search.substr(1).match(reg);
    if (r!=null) return unescape(r[2]); return null;
}

function generPageHtml(totalPage,totalCount){
    // var totalPage = {{ .totalPage }};
    // var totalRecords = {{.totalCount }};
    var totalPage = totalPage;
    var totalRecords = totalCount;
    var pageNo = getParameter('pno');
    var cate_id = getParameter('cate_id');
    if(!pageNo){
        pageNo = 1;
        //原本是实现无刷新跳转，我这是根据自己需求做的有刷新时跳转
        // 如：www.baidu.com/abcd/index.jhtml
        /*  let str=window.location.pathname;
            let two; // 第二个斜杠后内容
            let first = str.indexOf("/") + 1;
            let heng = str.indexOf("/", first);
            if (heng == -1) {
                } else {
                  two = str.substring(heng).substring(1, str.length);
                }
            if(two=="index.jhtml"){
                pageNo = 1;
            }else{
                pageNo = num;//num是根据自己要点击第几页写的
            }*/
    }
    //生成分页
    //有些参数是可选的，比如lang，若不传有默认值
    kkpager.generPageHtml({
        pno : pageNo,
        //总页码
        total : totalPage,
        //总数据条数
        totalRecords : totalRecords,
        mode : 'click',//默认值是link，可选link或者click
        click : function(n){
            this.selectPage(n);

            if(n==1){
                //原本是实现无刷新跳转，我这是根据自己需求做的有刷新时跳转
                //第一页写逻辑跳转
                // 如：window.location.href=......;
            }else{
                //除了第一页写逻辑跳转
            }
            if (cate_id){
                window.location.href='?cate_id='+cate_id+'&pno='+n
            }else{
                window.location.href='?'+'&pno='+n
            }

            return false;
        }

        /*      ,lang: {
                    firstPageText           : '首页',
                    firstPageTipText        : '首页',
                    lastPageText            : '尾页',
                    lastPageTipText         : '尾页',
                    prePageText             : '上一页',
                    prePageTipText          : '上一页',
                    nextPageText            : '下一页',
                    nextPageTipText         : '下一页',
                    totalPageBeforeText     : '共',
                    totalPageAfterText      : '页',
                    currPageBeforeText      : '当前第',
                    currPageAfterText       : '页',
                    totalInfoSplitStr       : '/',
                    totalRecordsBeforeText  : '共',
                    totalRecordsAfterText   : '条数据',
                    gopageBeforeText        : '&nbsp;转到',
                    gopageButtonOkText      : '确定',
                    gopageAfterText         : '页',
                    buttonTipBeforeText     : '第',
                    buttonTipAfterText      : '页'
                }*/
    });
}


window.onload = function(){
    if ($('#editor_content').length >0 ){
        init_editor("editor_content");
    }
}
// $(function () {
//     if ($('#editor_content').length >0 ){
//         init_editor("editor_content");
//     }
// })